/* Copyright (c) 2017
 *	Andrey Mikhalevich, Katren ltd.
 */
function BankList_Form(options){
	options = options || {};	
	
	options.formName = "BankList";
	options.controller = "Bank_Controller";
	options.method = "get_list";
	
	BankList_Form.superclass.constructor.call(this,options);
		
}
extend(BankList_Form,WindowFormObject);

/* Constants */


/* private members */

/* protected*/


/* public methods */

