/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2017
 
 * THIS FILE IS GENERATED FROM TEMPLATE build/templates/controllers/Controller_js20.xsl
 * ALL DIRECT MODIFICATIONS WILL BE LOST WITH THE NEXT BUILD PROCESS!!!
 
 * @class
 * @classdesc controller
 
 * @extends ControllerObjServer
  
 * @requires core/extend.js
 * @requires core/ControllerObjServer.js
  
 * @param {Object} options
 * @param {Model} options.listModelClass
 * @param {Model} options.objModelClass
 */ 

function ClientSearch_Controller(options){
	options = options || {};
	ClientSearch_Controller.superclass.constructor.call(this,options);	
	
	//methods
	this.add_search();
		
}
extend(ClientSearch_Controller,ControllerObjServer);

			ClientSearch_Controller.prototype.add_search = function(){
	var opts = {"controller":this};	
	var pm = new PublicMethodServer('search',opts);
	
				
	
	var options = {};
	
		options.required = true;
	
		options.maxlength = "250";
	
		pm.addField(new FieldString("query",options));
	
			
	this.addPublicMethod(pm);
}

		