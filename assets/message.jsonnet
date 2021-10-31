{%- with docs=message.ParentFile().SourceLocations().ByDescriptor(message) %}
{%- if message.LeadingComments %}
"#{{message.Name}}": d.obj("{{message.LeadingComments}}"),
{%- endif %}
{%- endwith -%}
{{message.Name}}:: {
	"#definition": d.obj("MessageDescriptor"),
	definition:: {{message|json|indent:1}},
	"#new": d.fn("Create {{message.FullName}}", [
		{% for i in range(message.Fields().Len()) %}
		{%- with field=message.Fields().Get(i) -%}
		d.arg({{field.JSONName|identifier|json}}, {% if field.IsList() -%}
			d.T.array
			{%- elif field.IsMap() -%}d.T.object
			{%- elif field.Kind == 8 -%}d.T.bool
			{%- elif field.Kind == 5 || field.Kind == 17 || field.Kind == 13 || field.Kind == 3 || field.Kind == 18 || field.Kind == 4 || field.Kind == 15 || field.Kind == 7 || field.Kind == 16 || field.Kind == 6 || field.Kind == 1 -%}d.T.int
			{%- elif field.Kind == 9 -%}d.T.string
			{%- elif field.Kind == 12 -%}d.T.string
			{%- elif field.Kind == 11 -%}d.T.object
			{%- elif field.Kind == 10 -%}d.T.object
			{%- else -%}
			d.T.any
			{%- endif -%}, {{field.Default|json}}), {% endwith %}
		{% endfor %}
	]),
	new(
		{%- for i in range(message.Fields().Len()) -%}
		{%- with field=message.Fields().Get(i) -%}
		{{field.JSONName|identifier}}={{field.Default|json}}
		{%- if not forloop.Last %}, {% endif -%}
		{%- endwith -%}
		{%- endfor -%}
	):: {
		_proto_type:: "{{message.FullName}}",
		{% for i in range(message.Fields().Len()) %}
		{%- with field=message.Fields().Get(i) -%}
		{{field.JSONName|json}}: {{field.JSONName|identifier}},
		{% if field.Kind == 11 && !field.IsMap() && !field.IsList() -%}{# struct #}
		assert self.{{field.JSONName}} == null || self.{{field.JSONName}}._proto_type == "{{field.Message.FullName}}" : "Value of '{{field.FullName}}' must be of type {{field.Message.FullName}}",
		{%- endif %}
		{% endwith -%}
		{%- endfor %}
	},

	{% for i in range(message.Fields().Len()) %}
	{%- with field=message.Fields().Get(i) -%}
	{%- with docs=field.ParentFile().SourceLocations().ByDescriptor(field) -%}
	{%- if docs.LeadingComments %}
	"#with{{field.JSONName|ucfirst}}": d.fn({{docs.LeadingComments|strip|json}}),
	{%- endif -%}
	{%- endwith %}
	with{{field.JSONName|ucfirst}}({{field.JSONName|identifier}}): { {{field.JSONName|json}}: {{field.JSONName|identifier}} },
	{% endwith -%}
	{%- endfor %}

	{% for i in range(message.Enums().Len()) %}
	{%- filter indent:1 -%}
	{%- include "enum.jsonnet" with enum=message.Enums().Get(i) -%}
	{%- endfilter %}
	{% endfor %}
	{% for i in range(message.Messages().Len()) %}
	{%- filter indent:1 -%}
	{%- include messageFile with message=message.Messages().Get(i) -%}
	{%- endfilter %}
	{% endfor %}
},
