package osbe 

const (
	ER_WRONG_TOKEN_FORMAT = "Неверный формат токена."
	ER_SESS_NOT_FOUND = "Сессия не нейдена."
	
	ER_PARSE_NO_FUNC = "Параметр 'func' не найден."
	ER_PARSE_NO_METH = "Не найдена функция в параметре."	
	ER_PARSE_NO_CONTR = "Не найден котроллер в параметре."	
	ER_PARSE_CTRL_NOT_DEFINED = "Контроллер '%s' не определен."
	
	ER_PARSE_NOT_VALID_EMPTY = "Поле: %s, пустое значение"
	
	ER_AUTH = "Ошибка авторизации."//100
	ER_AUTH_EXP = "Срок сессии истек."//101
	ER_AUTH_NOT_LOGGED = "Не авторизован."//102
	ER_AUTH_BANNED = "Доступ запрещен."//103
	ER_SQL_SERVER_CON = "Ошибка подключения к серверу базы данных."//105
	ER_SQL_QUERY = "Ошибка при выполнении запроса к базе данных."//106
	ER_VERSION = "Версии клиентского и серверного ПО отличается."//107
	ER_SESSION = "Ошибка работы с данными сессии."//109
	ER_COM_NO_CONTROLLER = "Контроллер не определен."//10
	ER_COM_METH_PROHIB = "Метод запрещен."//11
	ER_COM_NO_VIEW = "Вид не определен."//12
	ER_INTERNAL = "Server innner error."//13
	ER_DELETE_CONSTR_VIOL = "Удаление невозможно, так как существуют ссылки."//500
	ER_WRITE_CONSTR_VIOL = "Нарушение уникальности ключевых полей."//600
	ER_PM_INTERNAL = "Ошибка исполнения метода."//5
	
	ER_CONTOLLER_METH_NOT_DEFINED = "Метод '%s' контроллера '%s' не найден."
	
	ER_SQL_WHERE_FILED_CNT_MISMATCH = "Количество полей в условии не совпадает."
	ER_SQL_WHERE_UNKNOWN_COND = "Неизветсное условие '%s'."
	
	ER_UPDATE_EMPTY = "Поля для обновления не установлены."
	
	ER_NO_KEYS = "Ключевые поля не установлены."
	ER_NO_WHERE = "Не определено значение условия."
	
	ER_XSL_TRANSFORM = "Ошибка трансформации xsl шаблона."
	
)	

