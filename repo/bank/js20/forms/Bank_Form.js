/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2016
 
 * @class
 * @classdesc
 
 * @extends WindowFormObject	
 
 * @param {namespace} options
 * @param {string} options.formName
 * @param {string} options.method
 * @param {string} options.controller   
 */	
function Bank_Form(options){
	options = options || {};	
	
	options.formName = "Bank";
	options.method = "get_object";
	options.controller = "Bank_Controller";
	options.keys = {"bik":null};
	options.height = 500;
	options.width = 500;
	
	
	Bank_Form.superclass.constructor.call(this,options);
}
extend(Bank_Form,WindowFormObject);

/* Constants */


/* private members */

/* protected*/


/* public methods */

