/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2016
 
 * @class
 * @classdesc
	
 * @param {string} id view identifier
 * @param {namespace} options
 * @param {namespace} options.models All data models
 * @param {namespace} options.variantStorage {name,model}
 */	
function Bank_View(id,options){	

	options = options || {};
	
	options.enabled = false;
	
	options.cmdOk = false;
	options.cmdSave = false;
	
	Bank_View.superclass.constructor.call(this,id,options);
		
	this.addElement(new EditString(id+":bik",{
		"editMask":"999999999",
		"maxLenght":"9",
		"fixLength":true,
		"enabled":false,
		"cmdClear":false,
		"labelCaption":this.LB_CAP_BIK
	}));	

	this.addElement(new EditString(id+":name",{
		"maxLenght":"500",
		"enabled":false,
		"cmdClear":false,
		"labelCaption":this.LB_CAP_NAME
	}));	

	this.addElement(new EditString(id+":gor",{
		"maxLenght":"500",
		"enabled":false,
		"cmdClear":false,
		"labelCaption":this.LB_CAP_CITY
	}));	

	this.addElement(new EditString(id+":gr_descr",{
		"maxLenght":"50",
		"enabled":false,
		"cmdClear":false,
		"labelCaption":this.LB_CAP_GROUP
	}));	

	this.addElement(new EditString(id+":korshet",{
		"editMask":"99999999999999999999",
		"maxLenght":"20",
		"fixLength":true,
		"enabled":false,
		"cmdClear":false,
		"labelCaption":this.LB_CAP_BANK_ACC
	}));	
	this.addElement(new EditString(id+":adres",{
		"maxLenght":"500",
		"enabled":false,
		"cmdClear":false,
		"labelCaption":this.LB_CAP_ADDRESS
	}));	
	
	//****************************************************
	//read
	this.setReadPublicMethod((new Bank_Controller()).getPublicMethod("get_object"));
	this.m_model = options.models.BankList_Model;
	this.setDataBindings([
		new DataBinding({"control":this.getElement("bik"),"model":this.m_model}),
		new DataBinding({"control":this.getElement("name"),"model":this.m_model}),
		new DataBinding({"control":this.getElement("korshet"),"model":this.m_model}),
		new DataBinding({"control":this.getElement("gor"),"model":this.m_model}),
		new DataBinding({"control":this.getElement("gr_descr"),"model":this.m_model})
	]);	
}
extend(Bank_View,ViewObjectAjx);
