<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" 
 xmlns:html="http://www.w3.org/TR/REC-html40"
 xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
 xmlns:fo="http://www.w3.org/1999/XSL/Format">

<xsl:output method="xml"/> 

<xsl:variable name="CM_IN_PIX" select="0.02"/> 
<xsl:variable name="MIN_COL_WIDTH_PX" select="50"/> 
<!-- default widths for data types in px-->
<xsl:variable name="DEF_WIDTH_DATE" select="100"/>
<xsl:variable name="DEF_WIDTH_TIME" select="100"/>
<xsl:variable name="DEF_WIDTH_DATETIME" select="125"/>
<xsl:variable name="DEF_WIDTH_INT" select="65"/>
<xsl:variable name="DEF_WIDTH_DOUBLE" select="50"/>
<xsl:variable name="DEF_WIDTH_FILE" select="230"/>
<!--Data types-->
<xsl:variable name="DT_INT" select="'INTEGER'"/>
<xsl:variable name="DT_INT_UNSIGNED" select="'INTEGER_UNSIGNED'"/>
<xsl:variable name="DT_DOUBLE" select="'DOUBLE'"/>
<xsl:variable name="DT_DOUBLE_UNSIGNED" select="'DOUBLE_UNSIGNED'"/>
<xsl:variable name="DT_CUR_RUR" select="'CUR_RUR'"/>
<xsl:variable name="DT_CUR_USD" select="'CUR_USD'"/>
<xsl:variable name="DT_STRING" select="'VARCHAR'"/>
<xsl:variable name="DT_DATETIME" select="'DATETIME'"/>
<xsl:variable name="DT_DATE" select="'DATE'"/>
<xsl:variable name="DT_OBJECT" select="'OBJECT'"/>
<xsl:variable name="DT_FILE" select="'FILE'"/>
<xsl:variable name="DT_BOOL" select="'TINYINT'"/>
<xsl:variable name="DT_PWD" select="'PWD'"/>
<xsl:variable name="DT_EMAIL" select="'EMAIL'"/>
<xsl:variable name="DT_TEXT" select="'TEXT'"/>
<xsl:variable name="DT_ENUM" select="'ENUM'"/>
<xsl:variable name="DT_TIME" select="'TIME'"/>

<!-- Main template -->
<xsl:template match="/">
    <fo:root>
      <fo:layout-master-set>
        <fo:simple-page-master master-name="Report"
              page-height="29.7cm" page-width="21cm" margin-top="0.5cm"
			  margin-left="0.5cm" margin-right="0.5cm">
          <fo:region-body margin-top="0.3cm"/>
		  <fo:region-before/>
        </fo:simple-page-master>
      </fo:layout-master-set>
      <fo:page-sequence master-reference="Report">	  
	  
		<fo:static-content flow-name="xsl-region-before">
			<fo:block font-family="Arial" font-style="normal"
			font-weight="normal" text-align="left" font-size="8pt">
			Лист: <fo:page-number/>
			из <fo:page-number-citation ref-id="terminator"/>
		  </fo:block>
		</fo:static-content>
		<!-- flow-name="xsl-region-body"-->
        <fo:flow flow-name="xsl-region-body">
			<fo:block font-family="Arial" font-style="normal"
				font-weight="bold" text-align="center">
				<xsl:value-of select="page/repHeader/@title"/>
			</fo:block>
			
			<xsl:apply-templates select="document/model"/>
			
			<fo:block id="terminator"/>
        </fo:flow>			
		
      </fo:page-sequence>
    </fo:root>
</xsl:template>

<!-- table -->
<xsl:template match="model[not(@id='ModelServRespons')]">
	<fo:table table-layout="fixed">
	
		<!-- header-->
		<xsl:variable name="model_id" select="@id"/>
		<xsl:apply-templates select="/document/metadata[@modelId=$model_id]"/>
		
		<fo:table-body>
			<!-- all data rows -->
			<xsl:apply-templates select="row"/>
		</fo:table-body>
	</fo:table>
</xsl:template>

<!-- table header -->
<xsl:template match="metadata">
	<xsl:for-each select="field">
		<xsl:variable name="len">
			<xsl:choose>
				<xsl:when test="@displayWidth &gt; 0">
					<xsl:value-of select="@displayWidth"/>
				</xsl:when>
				<xsl:when test="@type=$DT_INT">
					<xsl:value-of select="$DEF_WIDTH_INT"/>
				</xsl:when>
				<xsl:when test="@type=$DT_DATE">
					<xsl:value-of select="$DEF_WIDTH_DATE"/>
				</xsl:when>
				<xsl:when test="@type=$DT_TIME">
					<xsl:value-of select="$DEF_WIDTH_TIME"/>
				</xsl:when>					
				<xsl:when test="@type=$DT_DATETIME">
					<xsl:value-of select="$DEF_WIDTH_DATETIME"/>
				</xsl:when>
				<xsl:when test="@type=$DT_FILE">
					<xsl:value-of select="$DEF_WIDTH_FILE"/>
				</xsl:when>
				<xsl:when test="@type=$DT_DOUBLE">
					<xsl:value-of select="$DEF_WIDTH_DOUBLE"/>
				</xsl:when>
				<xsl:otherwise>
					<xsl:value-of select="$MIN_COL_WIDTH_PX"/>
				</xsl:otherwise>
			</xsl:choose>
		</xsl:variable>	
		<xsl:variable name="len_cm" select="concat(round($CM_IN_PIX * $len),'cm')"/>
		<fo:table-column/>
		<!--column-width="{$len_cm}"-->
	</xsl:for-each>
	
	<fo:table-header text-align="center">
		<fo:table-row>
			<xsl:apply-templates select="field"/>
		</fo:table-row>
	</fo:table-header>
</xsl:template>

<!-- header columns -->
<xsl:template match="field">
	<xsl:variable name="col">
		<xsl:choose>
			<xsl:when test="@alias">
				<xsl:value-of select="@alias"/>
			</xsl:when>
			<xsl:otherwise test="@id">
				<xsl:value-of select="@id"/>
			</xsl:otherwise>			
		</xsl:choose>
	</xsl:variable>
	<fo:table-cell
		display-align="center"
		border-width="0.5mm" border-style="solid">
		<fo:block font-family="Arial" font-style="normal"
		font-weight="normal" text-align="center" font-size="8pt">
			<xsl:value-of select="$col"/>
		</fo:block>
	</fo:table-cell>
</xsl:template>

<!-- table row -->
<xsl:template match="row">
	<fo:table-row>
		<xsl:apply-templates/>
	</fo:table-row>
</xsl:template>

<!-- table cell -->
<xsl:template match="row/*">
	<xsl:variable name="field_num" select="position()"/>
	<fo:table-cell border-width="0.5mm" border-style="solid">						
		<fo:block font-family="Arial" font-style="normal"
		font-weight="normal" text-align="center" font-size="8pt">
			<xsl:value-of select="node()"/>
		</fo:block>
	</fo:table-cell>						
</xsl:template>

<!-- seconds transformation-->
<xsl:template name="echo_time_part">
	<xsl:param name="part"/>
	<xsl:choose>
		<xsl:when test="$part &gt; 9">
			<xsl:value-of select="$part"/>
		</xsl:when>
		<xsl:otherwise>
			0<xsl:value-of select="$part"/>
		</xsl:otherwise>
	</xsl:choose>
</xsl:template>

<xsl:template name="seconds_to_time">
	<xsl:param name="seconds"/>
	<xsl:variable name="days">
		<xsl:choose>
			<xsl:when test="$seconds/3600 &gt; 24">
				<xsl:value-of select="floor($seconds div 3600 div 24)"/>
			</xsl:when>
			<xsl:otherwise>
				<xsl:value-of select="0"/>
			</xsl:otherwise>
		</xsl:choose>
	</xsl:variable>
	<xsl:variable name="rest" select="($seconds - $days*24*3600) div 3600"/>
	<!--!!!seconds=<xsl:value-of select="$seconds"/>!!!-->
	<xsl:variable name="hours">
		<xsl:choose>
			<xsl:when test="$rest &gt; 0">
				<xsl:value-of select="floor($rest)"/>
			</xsl:when>
			<xsl:otherwise>
				<xsl:value-of select="0"/>
			</xsl:otherwise>
		</xsl:choose>
	</xsl:variable>
	<xsl:variable name="rest2" select="($seconds - $hours*3600 - $days*24*3600) div 60"/>
	<xsl:variable name="minutes">
		<xsl:choose>
			<xsl:when test="$rest2 &gt; 0">
				<xsl:value-of select="floor($rest2)"/>
			</xsl:when>
			<xsl:otherwise>
				<xsl:value-of select="0"/>
			</xsl:otherwise>
		</xsl:choose>
	</xsl:variable>
	
	<xsl:if test="$days &gt; 0">
		<xsl:value-of select="$days"/>.
	</xsl:if>
	<xsl:call-template name="echo_time_part">
		<xsl:with-param name="part" select="$hours"/>
	</xsl:call-template>:<xsl:call-template name="echo_time_part"><xsl:with-param name="part" select="$minutes"/>
	</xsl:call-template>	
	
</xsl:template>

</xsl:stylesheet>