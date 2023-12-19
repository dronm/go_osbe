/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2022
 
 * @extends ViewObjectAjx.js
 * @requires core/extend.js  
 * @requires controls/ViewObjectAjx.js 
 
 * @class
 * @classdesc
	
 * @param {string} id view identifier
 * @param {object} options
 * @param {object} options.models All data models
 * @param {object} options.variantStorage {name,model}
 */	
function PostDialog_View(id,options){	

	options = options || {};
	
	options.controller = new Post_Controller();
	options.model = (options.models&&options.models.Post_Model)? options.models.Post_Model : new Post_Model();
	
	options.addElement = function(){
		this.addElement(new EditString(id+":name",{
			"required": true,
			"focus": true,
			"maxLength":250,
			"labelCaption": "Наименование",
			"placeholder":"Наименование должности",
			"title":"Наименование должности"
		}));	

	}
	
	PostDialog_View.superclass.constructor.call(this,id,options);
	
	//****************************************************
	//read
	this.setDataBindings([
		new DataBinding({"control":this.getElement("name")})
	]);
	
	//write
	this.setWriteBindings([
		new CommandBinding({"control":this.getElement("name")})
	]);
}
extend(PostDialog_View,ViewObjectAjx);
