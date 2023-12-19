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

function Contact_Controller(options){
	options = options || {};
	options.listModelClass = ContactList_Model;
	options.objModelClass = ContactDialog_Model;
	Contact_Controller.superclass.constructor.call(this,options);	
	
	//methods
	this.addInsert();
	this.addUpdate();
	this.addDelete();
	this.addGetList();
	this.addGetObject();
	this.addComplete();
	this.add_upsert();
		
}
extend(Contact_Controller,ControllerObjServer);

			Contact_Controller.prototype.addInsert = function(){
	Contact_Controller.superclass.addInsert.call(this);
	
	var pm = this.getInsert();
	
	var options = {};
	options.primaryKey = true;options.autoInc = true;
	var field = new FieldInt("id",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Наименование";options.required = true;
	var field = new FieldString("name",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Должность";
	var field = new FieldInt("post_id",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Email";
	var field = new FieldString("email",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Телефон";
	var field = new FieldString("tel",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Добавочный номер";
	var field = new FieldString("tel_ext",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Описание для поиска";
	var field = new FieldText("descr",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Комментарий";
	var field = new FieldText("comment_text",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Адрес электр.почты подтвержден";
	var field = new FieldBool("email_confirmed",options);
	
	pm.addField(field);
	
	pm.addField(new FieldInt("ret_id",{}));
	
	
}

			Contact_Controller.prototype.addUpdate = function(){
	Contact_Controller.superclass.addUpdate.call(this);
	var pm = this.getUpdate();
	
	var options = {};
	options.primaryKey = true;options.autoInc = true;
	var field = new FieldInt("id",options);
	
	pm.addField(field);
	
	field = new FieldInt("old_id",{});
	pm.addField(field);
	
	var options = {};
	options.alias = "Наименование";
	var field = new FieldString("name",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Должность";
	var field = new FieldInt("post_id",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Email";
	var field = new FieldString("email",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Телефон";
	var field = new FieldString("tel",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Добавочный номер";
	var field = new FieldString("tel_ext",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Описание для поиска";
	var field = new FieldText("descr",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Комментарий";
	var field = new FieldText("comment_text",options);
	
	pm.addField(field);
	
	var options = {};
	options.alias = "Адрес электр.почты подтвержден";
	var field = new FieldBool("email_confirmed",options);
	
	pm.addField(field);
	
	
}

			Contact_Controller.prototype.addDelete = function(){
	Contact_Controller.superclass.addDelete.call(this);
	var pm = this.getDelete();
	var options = {"required":true};
		
	pm.addField(new FieldInt("id",options));
}

			Contact_Controller.prototype.addGetList = function(){
	Contact_Controller.superclass.addGetList.call(this);
	
	
	
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
	f_opts.alias = "Наименование";
	pm.addField(new FieldString("name",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldJSON("posts_ref",f_opts));
	var f_opts = {};
	f_opts.alias = "Email";
	pm.addField(new FieldString("email",f_opts));
	var f_opts = {};
	f_opts.alias = "Телефон";
	pm.addField(new FieldString("tel",f_opts));
	var f_opts = {};
	f_opts.alias = "Добавочный номер";
	pm.addField(new FieldString("tel_ext",f_opts));
	var f_opts = {};
	f_opts.alias = "Описание для поиска";
	pm.addField(new FieldText("descr",f_opts));
	var f_opts = {};
	f_opts.alias = "Комментарий";
	pm.addField(new FieldText("comment_text",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldText("tm_photo",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldString("tm_first_name",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldBool("tm_exists",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldBool("tm_activated",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldInt("ext_id",f_opts));
	pm.getField(this.PARAM_ORD_FIELDS).setValue("descr");
	
}

			Contact_Controller.prototype.addGetObject = function(){
	Contact_Controller.superclass.addGetObject.call(this);
	
	var pm = this.getGetObject();
	var f_opts = {};
		
	pm.addField(new FieldInt("id",f_opts));
	
	pm.addField(new FieldString("mode"));
	pm.addField(new FieldString("lsn"));
}

			Contact_Controller.prototype.addComplete = function(){
	Contact_Controller.superclass.addComplete.call(this);
	
	var f_opts = {};
	f_opts.alias = "";
	var pm = this.getComplete();
	pm.addField(new FieldText("descr",f_opts));
	pm.addField(new FieldInt("count", {}));
	pm.getField(this.PARAM_ORD_FIELDS).setValue("descr");	
}

			Contact_Controller.prototype.add_upsert = function(){
	var opts = {"controller":this};	
	var pm = new PublicMethodServer('upsert',opts);
	
				
	
	var options = {};
	
		options.required = true;
	
		options.maxlength = "250";
	
		pm.addField(new FieldString("name",options));
	
				
	
	var options = {};
	
		options.required = true;
	
		options.maxlength = "11";
	
		pm.addField(new FieldString("tel",options));
	
				
	
	var options = {};
	
		options.required = true;
	
		options.maxlength = "100";
	
		pm.addField(new FieldString("email",options));
	
				
	
	var options = {};
	
		options.required = true;
	
		options.maxlength = "20";
	
		pm.addField(new FieldString("tel_ext",options));
	
			
	this.addPublicMethod(pm);
}
			
		