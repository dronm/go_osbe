/** Copyright (c) 2019 
 *	Andrey Mikhalevich, Katren ltd.
 */
function PostEdit(id,options){

	options = options || {};	
	if (options.labelCaption!=""){
		options.labelCaption = options.labelCaption || "Должность:";
	}
	options.cmdInsert = (options.cmdInsert!=undefined)? options.cmdInsert:true;
	
	options.keyIds = options.keyIds || ["id"];
	
	//форма выбора из списка
	options.selectWinClass = PostList_Form;
	options.selectDescrIds = options.selectDescrIds || ["name"];
	
	//форма редактирования элемента
	options.editWinClass = Post_Form;
	
	options.acMinLengthForQuery = (options.acMinLengthForQuery!=undefined)? options.acMinLengthForQuery:1;
	options.acController = new Post_Controller();
	options.acPublicMethod = options.acController.getPublicMethod("complete");
	options.acModel = new Post_Model();
	options.acPatternFieldId = options.acPatternFieldId || "name";
	options.acKeyFields = options.acKeyFields || [options.acModel.getField("id")];
	options.acDescrFields = options.acDescrFields || [options.acModel.getField("name")];
	options.acICase = options.acICase || "1";
	options.acMid = options.acMid || "1";
	
	PostEdit.superclass.constructor.call(this,id,options);
}
extend(PostEdit,EditRef);

/* Constants */


/* private members */

/* protected*/


/* public methods */

