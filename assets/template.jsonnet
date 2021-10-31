local d = import "{{docUtils}}";

{
	package: {{input.Package|json}},

	{% for i in range(input.Enums().Len()) %}
	{%- filter indent:1 -%}
	{% include "enum.jsonnet" with enum=input.Enums().Get(i) %}
	{%- endfilter %}
	{% endfor %}

	{%- for i in range(input.Messages().Len()) %}
	{%- filter indent:1 -%}
	{% include "message.jsonnet" with message=input.Messages().Get(i) %}
	{%- endfilter %}
	{% endfor %}
}
