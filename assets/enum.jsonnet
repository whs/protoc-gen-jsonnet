{%- with docs=enum.ParentFile().SourceLocations().ByDescriptor(enum) %}
{%- if docs.LeadingComments %}
"#{{enum.Name}}": d.obj("{{docs.LeadingComments}}"),
{%- endif %}
{%- endwith -%}
{{enum.Name}}:: {
	"#_proto_def": d.val("EnumDescriptor"),
	_proto_def:: {{enum|json|indent:1}},
	"#_proto_type": d.val("Enum full type name"),
	_proto_type:: "{{enum.FullName}}",

	{% for i in range(enum.Values().Len()) %}
	{%- with value=enum.Values().Get(i) -%}
	{%- with docs=value.ParentFile().SourceLocations().ByDescriptor(value) -%}
	{%- if docs.LeadingComments -%}
	"#{{value.Number|identifier}}": d.val({{docs.LeadingComments|strip|json}}),
	"#{{value.Name|identifier}}": d.val({{docs.LeadingComments|strip|json}}),
	{%- endif %}
	"{{value.Number}}": {{value.Name|json}},
	{{value.Name|identifier}}: {{value.Name|json}},
	{% endwith -%}
	{% endwith -%}
	{%- endfor %}
},
