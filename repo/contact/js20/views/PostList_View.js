/** Copyright (c) 2022
 *	Andrey Mikhalevich, Katren ltd.
 *
 */
function PostList_View(id,options){	

	options = options || {};
	options.HEAD_TITLE = "Должности";

	PostList_View.superclass.constructor.call(this,id,options);
	
	var model = (options.models && options.models. Post_Model)? options.models. Post_Model : new Post_Model();
	var contr = new Post_Controller();
	
	var constants = {"doc_per_page_count":null,"grid_refresh_interval":null};
	window.getApp().getConstantManager().get(constants);
	
	var popup_menu = new PopUpMenu();
	var pagClass = window.getApp().getPaginationClass();
	
	this.addElement(new GridAjx(id+":grid",{
		"model":model,
		"controller":contr,
		"editInline":true,
		"editWinClass":null,
		"commands":new GridCmdContainerAjx(id+":grid:cmd"),		
		"popUpMenu":popup_menu,
		"filters":(options.detailFilters&&options.detailFilters.ConstructionSiteItemForManagerSupplyList_Model)? options.detailFilters.ConstructionSiteItemForManagerSupplyList_Model:null,
		"head":new GridHead(id+"-grid:head",{
			"elements":[
				new GridRow(id+":grid:head:row0",{
					"elements":[
						new GridCellHead(id+":grid:head:name",{
							"value":"Наименование",
							"columns":[
								new GridColumn({"field":model.getField("name")})
							],
							"sortable":true,
							"sort":"asc"							
						})
					]
				})
			]
		}),
		"pagination":new pagClass(id+"_page",
			{"countPerPage":constants.doc_per_page_count.getValue()}),		
		
		"autoRefresh":false,
		"refreshInterval":constants.grid_refresh_interval.getValue()*1000,
		"rowSelect":false,
		"focus":true
	}));	
	


}
extend(PostList_View,ViewAjxList);
