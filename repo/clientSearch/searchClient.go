package clientSearch

import(
	"context"
	"fmt"
	
	"dadata"	
	"osbe"
	"osbe/response"
	
	"github.com/jackc/pgx/v4"
)

var DadataKey string

const (
	RESP_SEARCH_NOT_FOUND_CODE = 1010
	RESP_SEARCH_NOT_FOUND_DESCR = "По запросу ничего не найдено"
	
	RESP_DATATA_NO_KEY_CODE = 1001
	RESP_DATATA_NO_KEY_DESCR = "Не задан ключ"
	
)

func searchClientByINN(app osbe.Applicationer, conn *pgx.Conn, query string) (*dadata.FindByIdResponse, error) {
	search_res := &dadata.FindByIdResponse{}
	row := conn.QueryRow(context.Background(),
		`SELECT
			response
		FROM client_search.dadata_cache
		WHERE query = $1`,
		query,
	)
	err := row.Scan(search_res)
	if err != nil && err != pgx.ErrNoRows {
		return nil, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("searchClientByINN conn.QueryRow(): %v",err))
		
	}else if err != nil && err == pgx.ErrNoRows {
		//dadata search
		
		//dadata key
		if DadataKey == "" {
			return nil, osbe.NewPublicMethodError(RESP_DATATA_NO_KEY_CODE, RESP_DATATA_NO_KEY_DESCR)
		}
		
		dadata_api := dadata.NewDadata(DadataKey)
		search_res, err = dadata_api.FindById(query, "", dadata.BY_ID_SEARCH_BRANCH_TYPE_ALL, dadata.BY_ID_SEARCH_TYPE_ALL, 1)
		if err != nil {
			return nil, osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("searchClientByINN dadata_api.FindById(): %v",err))
		}
		
		if len(search_res.Suggestions) == 0 || search_res.Suggestions[0].Data == nil  {
			return nil, osbe.NewPublicMethodError(RESP_SEARCH_NOT_FOUND_CODE, RESP_SEARCH_NOT_FOUND_DESCR)
		}
		
		//ON CONFLICT (query) DO UPDATE SET response = $2`,
		//does not work on foreign table	
		if _, err := conn.Exec(context.Background(),
			`INSERT INTO client_search.dadata_cache (query, response)
			VALUES ($1, $2)`,		
			query,
			search_res,
		); err != nil {
			//cache error, still return result
			app.GetLogger().Errorf("searchClientByINN conn.Exec() INSERT failed: %v",err)
			return search_res, nil
		}
		
	}
	
	return search_res, nil
}
/*
func fillClientByINN(app osbe.Applicationer, conn *pgx.Conn, inn string, clientRow *ClientDialog) error {
	//organization does not exists!
	search_res, err := searchClientByINN(app, conn, inn)
	if err != nil {
		//dadata error
		return err
	}
	
	data := search_res.Suggestions[0].Data	
	if *data.Type == dadata.ORG_TYPE_ENTERPRIZE && data.Name != nil && data.Name.Short != nil && *data.Name.Short !="" {
		clientRow.Name.SetValue(*data.Name.Short)
		
	}else if *data.Type == dadata.ORG_TYPE_PERSON && data.Fio != nil && data.Fio.Surname != nil && *data.Fio.Surname !="" {
		n := *data.Fio.Surname
		if data.Fio.Name != nil && *data.Fio.Name !="" {
			n+= " "+*data.Fio.Name
		}
		if data.Fio.Patronymic != nil && *data.Fio.Patronymic !="" {
			n+= " "+*data.Fio.Patronymic
		}
		clientRow.Name.SetValue(n)
		
	}else{
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL,"Client_Controller_check_for_register org name is null")	
	}	
	clientRow.Inn.SetValue(inn)
	
	if data.Name != nil && data.Name.Full_with_opf != nil {
		clientRow.Name_full.SetValue(*data.Name.Full_with_opf)
	}
	
	//if data.Management != nil && data.Management.Name != nil {
	//}
	
	//if data.Management != nil && data.Management.Post != nil {
	//}
	if data.Kpp != nil {
		clientRow.Kpp.SetValue(*data.Kpp)
	}
	if data.Ogrn != nil {
		clientRow.Ogrn.SetValue(*data.Ogrn)
	}
	if data.Okpo != nil {
		clientRow.Okpo.SetValue(*data.Okpo)
	}
	if data.Okved != nil {
		clientRow.Okved.SetValue(*data.Okved)
	}
	if data.Address != nil && data.Address.Unrestricted_value != nil {
		clientRow.Legal_address.SetValue(*data.Address.Unrestricted_value)
	}
	
	if data.Phones != nil && len(data.Phones) > 0 {
		tels := ""
		for _, ph := range data.Phones {
			if tels != "" {
				tels += ", "
			}
			tels+= *ph.Source
		}
		clientRow.Tel.SetValue(tels)
	}
	if data.Emails != nil && len(data.Emails) > 0 {
		emails := ""
		for _, e := range data.Emails {
			emails = *e.Source
			break
		}
		clientRow.Email.SetValue(emails)
	}
	
	return nil
}
*/
