/** Copyright (c) 2019
 *	Andrey Mikhalevich, Katren ltd.
 */
function ContactList_Form(options){
	options = options || {};	
	
	options.formName = "ContactList";
	options.controller = "Contact_Controller";
	options.method = "get_list";
	
	ContactList_Form.superclass.constructor.call(this,options);
		
}
extend(ContactList_Form,WindowFormObject);

/* Constants */


/* private members */

/* protected*/


/* public methods */

