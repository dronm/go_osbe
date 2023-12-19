package userOperation

import(
	"context"
	"fmt"
	
	"osbe/logger"
	
	"github.com/jackc/pgx/v4"
)

func StartUserOperation(conn *pgx.Conn, operation string, userID int64, operationID string) error {
	if _, err := conn.Exec(context.Background(), 
		`INSERT INTO user_operations (user_id, operation_id, operation, status) VALUES ($1, $2, $3, 'start')
		ON CONFLICT (user_id, operation_id) DO UPDATE
		SET
			user_id = $1,
			operation = $3,
			status = 'start'`,
		userID,
		operationID,
		operation,
	); err != nil {
		return err
	}
	return nil		
}

func EndUserOperationWithError(log logger.Logger, conn *pgx.Conn, userID int64, operationID string, userErr error) {
	er_descr := ""
	if userErr == nil {
		er_descr = "Неизвестная ошибка"
	}else{
		er_descr = userErr.Error()
	}
	if _, err := conn.Exec(context.Background(), 
		`UPDATE user_operations
		SET
			status = 'end',
			error_text = $1,
			date_time_end = now()
		WHERE user_id = $2 AND operation_id = $3`,
		er_descr,
		userID,
		operationID,				
	); err != nil {
		log.Errorf("EndUserOperationWithError conn.Exec(): %v", err)
	}		
}

func SetUserOperationComment(log logger.Logger, conn *pgx.Conn, userID int64, operationID string, commentText string) {
	if _, err := conn.Exec(context.Background(), 
		`UPDATE user_operations
		SET
			comment_text = CASE WHEN comment_text IS NULL THEN $1 ELSE comment_text||', '||$1 END
		WHERE user_id = $2 AND operation_id = $3`,
		commentText,
		userID,
		operationID,				
	); err != nil {
		log.Errorf("setUserOperationComment conn.Exec(): %v", err)
	}		
}

func EndUserOperation(log logger.Logger, conn *pgx.Conn, userID int64, operationID string) {
	if _, err := conn.Exec(context.Background(), 
		`UPDATE user_operations
		SET
			status = 'end',
			date_time_end = now(),
			end_wal_lsn = pg_current_wal_lsn()::text
		WHERE user_id = $1 AND operation_id = $2`,
		userID,
		operationID,
	); err != nil {
		EndUserOperationWithError(log, conn, userID, operationID, fmt.Errorf("EndUserOperation UPDATE user_operations: %v", err))
	}
}

func ProgressEvent(conn *pgx.Conn, operationID, fieldId, comment string, res bool) error {
	res_s := "true"
	if !res {
		res_s = "false"
	}
	if _, err := conn.Exec(context.Background(), fmt.Sprintf(`SELECT pg_notify('UserOperation.%s', json_build_object('params', json_build_object('status', 'progress', 'f', '%s', 'c', '%s', 'res', %s))::text)`,
		operationID, fieldId, comment, res_s)); err != nil {
		return err
	}
	return nil
}

