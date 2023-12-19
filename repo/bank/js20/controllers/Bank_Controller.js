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

function Bank_Controller(options){
	options = options || {};
	options.listModelClass = BankList_Model;
	options.objModelClass = BankList_Model;
	Bank_Controller.superclass.constructor.call(this,options);	
	
	//methods
	this.addGetList();
	this.addGetObject();
	this.addComplete();
		
}
extend(Bank_Controller,ControllerObjServer);

			Bank_Controller.prototype.addGetList = function(){
	Bank_Controller.superclass.addGetList.call(this);
	
	
	
	var pm = this.getGetList();
	
	pm.addField(new FieldInt(this.PARAM_COUNT));
	pm.addField(new FieldInt(this.PARAM_FROM));
	pm.addField(new FieldString(this.PARAM_COND_FIELDS));
	pm.addField(new FieldString(this.PARAM_COND_SGNS));
	pm.addField(new FieldString(this.PARAM_COND_VALS));
	pm.addField(new FieldString(this.PARAM_COND_ICASE));
	pm.addField(new FieldString(this.PARAM_ORD_FIELDS));
	pm.addField(new FieldString(this.PARAM_ORD_DIRECTS));
	pm.addField(new FieldString(this.PARAM_FIELD_SEP));
	pm.addField(new FieldString(this.PARAM_FIELD_LSN));

	var f_opts = {};
	f_opts.alias = "БИК";
	pm.addField(new FieldString("bik",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldString("codegr",f_opts));
	var f_opts = {};
	f_opts.alias = "Регион";
	pm.addField(new FieldString("gr_descr",f_opts));
	var f_opts = {};
	f_opts.alias = "Наименование";
	pm.addField(new FieldText("name",f_opts));
	var f_opts = {};
	f_opts.alias = "Кoр.счет";
	pm.addField(new FieldString("korshet",f_opts));
	var f_opts = {};
	f_opts.alias = "Адрес";
	pm.addField(new FieldText("adres",f_opts));
	var f_opts = {};
	f_opts.alias = "Город";
	pm.addField(new FieldText("gor",f_opts));
	var f_opts = {};
	
	pm.addField(new FieldInt("tgoup",f_opts));
	pm.getField(this.PARAM_ORD_FIELDS).setValue("codegr,bik");
	
}

			Bank_Controller.prototype.addGetObject = function(){
	Bank_Controller.superclass.addGetObject.call(this);
	
	var pm = this.getGetObject();
	var f_opts = {};
	f_opts.alias = "БИК";	
	pm.addField(new FieldString("bik",f_opts));
	
	pm.addField(new FieldString("mode"));
	pm.addField(new FieldString("lsn"));
}

			Bank_Controller.prototype.addComplete = function(){
	Bank_Controller.superclass.addComplete.call(this);
	
	var f_opts = {};
	f_opts.alias = "";
	var pm = this.getComplete();
	pm.addField(new FieldString("bik",f_opts));
	pm.getField(this.PARAM_ORD_FIELDS).setValue("bik");	
}

		