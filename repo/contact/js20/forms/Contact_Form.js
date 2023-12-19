/** Copyright (c) 2022
	Andrey Mikhalevich, Katren ltd.
*/
function Contact_Form(options){
	options = options || {};	
	
	options.formName = "ContactDialog";
	options.controller = "Contact_Controller";
	options.method = "get_object";
	
	Contact_Form.superclass.constructor.call(this,options);
	
}
extend(Contact_Form,WindowFormObject);

/* Constants */


/* private members */

/* protected*/


/* public methods */

