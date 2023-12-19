/** Copyright (c) 2022
 *	Andrey Mikhalevich, Katren ltd.
 */
function Post_Form(options){
	options = options || {};	
	
	options.formName = "PostDialog";
	options.controller = "Post_Controller";
	options.method = "get_object";
	
	Post_Form.superclass.constructor.call(this,options);
	
}
extend(Post_Form,WindowFormObject);

/* Constants */


/* private members */

/* protected*/


/* public methods */

