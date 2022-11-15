<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:output method="html" indent="yes"
			doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN" 
			doctype-system="http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"/>

<xsl:variable name="BASE_PATH" select="/document/model[@id='ModelVars']/row[1]/basePath"/>
<xsl:variable name="VERSION" select="/document/model[@id='ModelVars']/row[1]/scriptId"/>
<xsl:variable name="TITLE" select="/document/model[@id='ModelVars']/row[1]/title"/>

<xsl:template match="/">
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=UTF-8"/>
		<xsl:apply-templates select="/document/model[@id='ModelVars']"/>
		<xsl:apply-templates select="/document/model[@id='ModelStyleSheet']/row"/>
		<script>
			var HOST_NAME = '<xsl:value-of select="/document/model[@id='ModelVars']/row/basePath"/>';
			var Connect;
			var MainView;//current opened view
			var onViewClose;//
			var CONSTANTS;
			
			function pageLoad(){
				CONSTANTS = new Constant_Controller(new ServConnector(HOST_NAME));
				var v = new <xsl:value-of select="/document/model[@id='ModelVars']/row/jsViewId"/>("MainView",
				{connect:new ServConnector('<xsl:value-of select="/document/model[@id='ModelVars']/row/basePath"/>'),
				onClose:function(){
					window.close();
				}
				});
				v.toDOM(document.body);				
				//document.head.title.innerHTML=
				//v.getTitleControl().getValue();
			}
		</script>		
		<title><xsl:value-of select="$TITLE"/></title>
	</head>
	<body onload="pageLoad();">
		<!--waiting  -->
		<div id="waiting">
			<div>Загрузка библиотек...</div>
			<img src="{$BASE_PATH}img/common/wait.gif" alt="загрузка"/>
		</div>
		
		<!--ALL js modules -->
		<xsl:apply-templates select="/document/model[@id='ModelJavaScript']/row"/>
		<script>
			var dv = document.getElementById("waiting");
			if (dv!==null){
				dv.parentNode.removeChild(dv);
			}
		</script>		
	</body>
</html>		
</xsl:template>

<xsl:template match="model[@id='ModelVars']/row">
	<xsl:if test="author">
		<meta name="Author" content="{author}"></meta>
	</xsl:if>
	<xsl:if test="keywords">
		<meta name="Keywords" content="{keywords}"></meta>
	</xsl:if>
	<xsl:if test="description">
		<meta name="Description" content="{description}"></meta>
	</xsl:if>
	
</xsl:template>

<xsl:template match="model[@id='ModelStyleSheet']/row">
	<link rel="stylesheet" href="{concat($BASE_PATH,href,'?',$VERSION)}" type="text/css"/>
</xsl:template>

<xsl:template match="model[@id='ModelJavaScript']/row">
	<script src="{concat($BASE_PATH,href,'?',$VERSION)}"></script>
</xsl:template>

<xsl:template match="model[@id='ModelServResponse']/row">
	<xsl:if test="result/node()='1'">
	<div class="error"><xsl:value-of select="descr"/></div>
	</xsl:if>
</xsl:template>

</xsl:stylesheet>