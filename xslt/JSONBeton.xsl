<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:import href="JSON.xsl"/>

<xsl:template match="menuItem">
	<xsl:variable name="level" select="count(ancestor::*)-2"/>
	<xsl:if test="(count(preceding-sibling::*) &gt; 0) or ($level &gt; 0)">,</xsl:if>	
	{"id":"<xsl:value-of select="generate-id()"/>",
	"level":"<xsl:value-of select="$level"/>",
	"parentId":"<xsl:if test="$level &gt;0"><xsl:value-of select="generate-id(..)"/></xsl:if>",
	"descr":"<xsl:value-of select="@descr"/>","viewId":"<xsl:value-of select="@viewId"/>","default":"<xsl:value-of select="@default"/>"
	}
	<xsl:apply-templates/>	
</xsl:template>

</xsl:stylesheet>