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

function EntityContact_Controller(options){
	options = options || {};
	options.listModelClass = EntityContactList_Model;
	options.objModelClass = EntityContactList_Model;
	EntityContact_Controller.superclass.constructor.call(this,options);	
	
	//methods
	this.addInsert();
	this.addUpdate();
	this.addDelete();
	this.addGetList();
	this.addGetObject();
		
}
extend(EntityContact_Controller,ControllerObjServer);

			EntityContact_Controller.prototype.addInsert = function(){
	EntityContact_Controller.superclass.addInsert.call(this);
	
	var pm = this.getInsert();
	
	var options = {};
	options.primaryKey = true;options.autoInc = true;
	var field = new FieldInt("id",options);
	
	pm.addField(field);
	
	var options = {};
	options.required = true;	
	options.enumValues = 'users';
	var field = new FieldEnum("entity_type",options);
	
	pm.addField(field);
	
	var options = {};
	options.required = true;
	var field = new FieldInt("entity_id",options);
	
	pm.addField(field);
	
	var options = {};
	options.required = true;
	var field = new FieldInt("contact_id",options);
	
	pm.addField(field);
	
	var options = {};
	options.required = true;
	var field = new FieldDateTimeTZ("mod_date_time",options);
	
	pm.addField(field);
	
	pm.addField(new FieldInt("ret_id",{}));
	
	
}

			EntityContact_Controller.prototype.addUpdate = function(){
	EntityContact_Controller.superclass.addUpdate.call(this);
	var pm = this.getUpdate();
	
	var options = {};
	options.primaryKey = true;options.autoInc = true;
	var field = new FieldInt("id",options);
	
	pm.addField(field);
	
	field = new FieldInt("old_id",{});
	pm.addField(field);
	
	var options = {};
		
	options.enumValues = 'users';
	options.enumValues+= (options.enumValues=='')? '':',';
	options.enumValues+= 'null';
	
	var field = new FieldEnum("entity_type",options);
	
	pm.addField(field);
	
	var options = {};
	
	var field = new FieldInt("entity_id",options);
	
	pm.addField(field);
	
	var options = {};
	
	var field = new FieldInt("contact_id",options);
	
	pm.addField(field);
	
	var options = {};
	
	var field = new FieldDateTimeTZ("mod_date_time",options);
	
	pm.addField(field);
	
	
}

			EntityContact_Controller.prototype.addDelete = function(){
	EntityContact_Controller.superclass.addDelete.call(this);
	var pm = this.getDelete();
	var options = {"required":true};
		
	pm.addField(new FieldInt("id",options));
}

			EntityContact_Controller.prototype.addGetList = function(){
	EntityContact_Controller.superclass.addGetList.call(this);
	
	
	
	var pm = this.getGetList();
	
	pm.addField(new FieldInt(this.PARAM_COUNT));
	pm.addField(new FieldInt(this.PARAM_FROM));
	pm.addField(new FieldString(this.PARAM_COND_FIELDS));
	pm.addField(new FieldString(this.PARAM_COND_SGNS));
	pm.addField(new FieldString(this.PARAM_COND_VALS));
	pm.addField(new FieldString(this.PARAM_COND_JOINS));
	pm.addField(new FieldString(this.PARAM_COND_ICASE));
	pm.addField(new FieldString(this.PARAM_ORD_FIELDS));
	pm.addField(new FieldString(this.PARAM_ORD_DIRECTS));
	pm.addField(new FieldString(this.PARAM_FIELD_SEP));
	pm.addField(new FieldString(this.PARAM_FIELD_LSN));
	pm.addField(new FieldString(this.PARAM_EXP_FNAME));

	var f_opts = {};
	
	pm.addField(new FieldInt("id",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldString("entity_type",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldInt("entity_id",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldJSON("entities_ref",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldInt("contact_id",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldString("contact_name",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldJSON("contact_attrs",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldJSON("contacts_ref",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldBool("tm_exists",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldBool("tm_activated",f_opts));
}

			EntityContact_Controller.prototype.addGetObject = function(){
	EntityContact_Controller.superclass.addGetObject.call(this);
	
	var pm = this.getGetObject();
	var f_opts = {};
		
	pm.addField(new FieldInt("id",f_opts));
	
	pm.addField(new FieldString("mode"));
	pm.addField(new FieldString("lsn"));
}

		