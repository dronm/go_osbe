/*
 * Copyright (c) 2022
 * Andrey Mikhalevich, Katren ltd.
 */
function PostList_Form(options){
	options = options || {};	
	
	options.formName = "PostList";
	options.controller = "Post_Controller";
	options.method = "get_list";
	
	PostList_Form.superclass.constructor.call(this,options);
		
}
extend(PostList_Form,WindowFormObject);

/* Constants */


/* private members */

/* protected*/


/* public methods */

