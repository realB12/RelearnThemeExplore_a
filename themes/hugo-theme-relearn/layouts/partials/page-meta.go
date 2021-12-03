<!--
  Attention: Pls. note that syntactical errors in this file might NOT lead to errors during Hugo
             generation!

  This page-meta.go template is called for every single page (list and single) when being rendered
  as part of the following rendering hierarchy:

    header.html
    |
    +- meta.html:
       |
       +- page-meta.go

  This partial just sets some internal Scratchpad-variables about the page's context such as pointers to
  1. next page
  2. previous page
  3. sub pages
  4. whether it is hidden or not
  etc.

  However this page-meta template does NOT RENDER anything on screen! It just delivers the page-scoped
  Scratchpad-variables required for the remaining templates to get the stuff (e.g. arrows for previous
  and next pages) properly rendered.
  ---

  Notes about the Scratch functions:

  The following makes use of the Hugo Provided Scratch framework.
  Find the docs about scratching here:
  ->  https://www.regisphilibert.com/blog/2017/04/hugo-scratch-explained-variable/

-->

{{- $currentNode := . }}

<!-- initially deleting all scratch-variables for this node to make sure, that no old values are used -->
{{- $currentNode.Scratch.Delete "relearnIsSelfFound"  }}
{{- $currentNode.Scratch.Delete "relearnPrevPage"     }}
{{- $currentNode.Scratch.Delete "relearnNextPage"     }}
{{- $currentNode.Scratch.Delete "relearnSubPages"     }}
{{- $currentNode.Scratch.Delete "relearnIsHiddenNode" }}
{{- $currentNode.Scratch.Delete "relearnIsHiddenStem" }}
{{- $currentNode.Scratch.Delete "relearnIsHiddenFrom" }}


<!-- calling the template-function with the Site.Home as the starting node
     Every call is called recursively with 4 parameters:
     1. node
     2. currentNode
     3. hiddenstem
     4. hiddencurrent
-->
{{- template "relearn-structure" dict "node" .Site.Home "currentnode" $currentNode "hiddenstem" false "hiddencurrent" false }}

<!-- implementing the above called RECURSIVELY called template() "relearn-structure" function
     that will be called ONCE with the Homepage (Site.Hompe)
-->
{{- define "relearn-structure" }}

	{{- $currentNode := .currentnode }}

	{{- $isSelf := eq $currentNode.RelPermalink .node.RelPermalink }}

	{{- $isDescendant := and (not $isSelf) (.node.IsDescendant $currentNode) }}

	{{- $isAncestor := and (not $isSelf) (.node.IsAncestor $currentNode) }}

	{{- $isOther := and (not $isDescendant) (not $isSelf) (not $isAncestor) }}

	{{- if $isSelf }}
		{{- $currentNode.Scratch.Set "relearnIsSelfFound" true }}
	{{- end}}

	{{- $isSelfFound := eq ($currentNode.Scratch.Get "relearnIsSelfFound") true }}

  {{- $isPreSelf := and (not $isSelfFound) (not $isSelf) }}

  {{- $isPostSelf := and ($isSelfFound) (not $isSelf) }}

	{{- $hidden_node := or (.node.Params.hidden) (eq .node.Title "") }}
	{{- $hidden_stem:= or $hidden_node .hiddenstem }}
	{{- $hidden_current_stem:= or $hidden_node .hiddencurrent }}
	{{- $hidden_from_current := or (and $hidden_node (not $isAncestor) (not $isSelf) ) (and .hiddencurrent (or $isPreSelf $isPostSelf $isDescendant) ) }}

	{{- $currentNode.Scratch.SetInMap "relearnIsHiddenNode" .node.RelPermalink $hidden_node }}
	{{- $currentNode.Scratch.SetInMap "relearnIsHiddenStem" .node.RelPermalink $hidden_stem }}
	{{- $currentNode.Scratch.SetInMap "relearnIsHiddenFrom" .node.RelPermalink $hidden_current_stem }}

	{{- if not $hidden_from_current }}
		{{- if $isPreSelf }}
			{{- $currentNode.Scratch.Set "relearnPrevPage" .node }}
		{{- else if and $isPostSelf (eq ($currentNode.Scratch.Get "relearnNextPage") nil) }}
			{{- $currentNode.Scratch.Set "relearnNextPage" .node }}
		{{- end}}
	{{- end }}

	{{- $currentNode.Scratch.Set "relearnSubPages" .node.Pages }}

	{{- if .node.IsHome }}
		{{- $currentNode.Scratch.Set "relearnSubPages" .node.Sections }}
	{{- else if .node.Sections }}
		{{- $currentNode.Scratch.Set "relearnSubPages" (.node.Pages | union .node.Sections) }}
	{{- end }}

  <!-- pages contains a point to all subpages of this node -->
	{{- $pages := ($currentNode.Scratch.Get "relearnSubPages") }}

	{{- $defaultOrdersectionsby := .Site.Params.ordersectionsby | default "weight" }}
	{{- $currentOrdersectionsby := .node.Params.ordersectionsby | default $defaultOrdersectionsby }}

	{{- if eq $currentOrdersectionsby "title"}}
		{{- range $pages.ByTitle  }}
			{{- template "relearn-structure" dict "node" . "currentnode" $currentNode "hiddenstem" $hidden_stem "hiddencurrent" $hidden_from_current }}
		{{- end}}
	{{- else}}
		{{- range $pages.ByWeight  }}
			{{- template "relearn-structure" dict "node" . "currentnode" $currentNode "hiddenstem" $hidden_stem "hiddencurrent" $hidden_from_current }}
		{{- end}}
	{{- end }}

{{- end }}

