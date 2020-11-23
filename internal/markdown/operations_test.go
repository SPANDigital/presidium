package markdown

import (
	"github.com/magiconair/properties/assert"
	"testing"
)


func TestCamelCase(t *testing.T) {
	input := `apple_scope`
	expectedOutput := `appleScope`
	output := ensureCamelCase(input)
	assert.Equal(t, output, expectedOutput)

	input = `wow_that_is_a_long_snake`
	expectedOutput = `wowThatIsALongSnake`
	output = ensureCamelCase(input)
	assert.Equal(t, output, expectedOutput)

	input = `CamelCase`
	expectedOutput = `CamelCase`
	output = ensureCamelCase(input)
	assert.Equal(t, output, expectedOutput)

	input = `camelCase`
	expectedOutput = `camelCase`
	output = ensureCamelCase(input)
	assert.Equal(t, output, expectedOutput)
}


func TestSimpleStatement(t *testing.T) {
	input := `site.audience == "Cinemo"`
	expectedOutput := `(eq $.Page.Site.Params.audience "Cinemo")`
	output := parseIfCondition(input)
	assert.Equal(t, output, expectedOutput)
}

func TestSimpleStatement__Shortcode(t *testing.T) {
	input := `site.audience == "Cinemo"`
	expectedOutput := `"audience" "Cinemo"`
	output := parseIfConditionShortcode(input)
	assert.Equal(t, output, expectedOutput)
}

func TestNeqStatement(t *testing.T) {
	input := `site.audience != "Cinemo"`
	expectedOutput := `(ne $.Page.Site.Params.audience "Cinemo")`
	output := parseIfCondition(input)
	assert.Equal(t, output, expectedOutput)
}

func TestCompoundStatement(t *testing.T) {
	input := `site.apple_scope == "internal" and site.audience != "Cinemo"`
	expectedOutput := `(eq $.Page.Site.Params.appleScope "internal") | and (ne $.Page.Site.Params.audience "Cinemo")`
	output := parseIfCondition(input)
	assert.Equal(t, output, expectedOutput)
}

func TestCompoundStatement__Shortcode(t *testing.T) {
	input := `site.apple_scope == "internal" and site.audience == "Cinemo"`
	expectedOutput := `"appleScope" "internal" "audience" "Cinemo"`
	output := parseIfConditionShortcode(input)
	assert.Equal(t, output, expectedOutput)
}

func TestFullIfStatement(t *testing.T) {
	input := `{% if site.audience == "Cinemo" or site.scope != "internal" %} Some Text {% elsif site.audience == "MusicKitJS" %} More Text {% else %} Default Text {% endif %}`
	expectedOutput := `{{% if.inline %}}{{ if (eq $.Page.Site.Params.audience "Cinemo") | or (ne $.Page.Site.Params.scope "internal") }} Some Text {{ else if (eq $.Page.Site.Params.audience "MusicKitJS") }} More Text {{ else }} Default Text {{ end }}{{% /if.inline %}}`
	output := parseIfStatements(input)
	assert.Equal(t, output, expectedOutput)
}

func TestFullIfStatement__Shortcode(t *testing.T) {
	input := `{% if site.audience == "Cinemo" or site.scope == "internal" %} {{% innershortcode %}} {% elsif site.audience == "MusicKitJS" %} More Text {% else %} Default Text {% endif %}`
	expectedOutput := `{{% when "audience" "Cinemo" "scope" "internal" %}} {{% innershortcode %}} {{% /when %}}{{% when "audience" "MusicKitJS" %}} More Text {{% /when %}}{{% default %}} Default Text {{% /default %}}`
	output := parseIfStatements(input)
	assert.Equal(t, output, expectedOutput)
}

func TestMultipleLinesOfIfStatements(t *testing.T) {
	input := `...
{% if site.audience == "Cinemo" or site.scope != "internal" %}
Some Text
{% elsif site.audience == "MusicKitJS" %}
More Text
{% else %}
Default Text
{% endif %}
...
{% if site.audience == "Cinemo" or site.scope != "internal" %}
Some Text
{% elsif site.audience == "MusicKitJS" %}
More Text
{% else %}
Default Text
{% endif %}
...`
	expectedOutput := `...
{{% if.inline %}}{{ if (eq $.Page.Site.Params.audience "Cinemo") | or (ne $.Page.Site.Params.scope "internal") }}
Some Text
{{ else if (eq $.Page.Site.Params.audience "MusicKitJS") }}
More Text
{{ else }}
Default Text
{{ end }}{{% /if.inline %}}
...
{{% if.inline %}}{{ if (eq $.Page.Site.Params.audience "Cinemo") | or (ne $.Page.Site.Params.scope "internal") }}
Some Text
{{ else if (eq $.Page.Site.Params.audience "MusicKitJS") }}
More Text
{{ else }}
Default Text
{{ end }}{{% /if.inline %}}
...`
	output := parseIfStatements(input)
	assert.Equal(t, output, expectedOutput)
}

func TestMultipleLinesOfIfStatements__Shortcode(t *testing.T) {
	input := `...
{% if site.audience == "Cinemo" and site.scope == "internal" %}
{{% someshortcode %}}
{% elsif site.audience == "MusicKitJS" %}
More Text
{% else %}
Default Text
{% endif %}
...
{% if site.audience == "Cinemo" or site.scope != "internal" %}
Some Text
{% elsif site.audience == "MusicKitJS" %}
More Text
{% else %}
Default Text
{% endif %}
...`
	expectedOutput := `...
{{% when "audience" "Cinemo" "scope" "internal" %}}
{{% someshortcode %}}
{{% /when %}}{{% when "audience" "MusicKitJS" %}}
More Text
{{% /when %}}{{% default %}}
Default Text
{{% /default %}}
...
{{% if.inline %}}{{ if (eq $.Page.Site.Params.audience "Cinemo") | or (ne $.Page.Site.Params.scope "internal") }}
Some Text
{{ else if (eq $.Page.Site.Params.audience "MusicKitJS") }}
More Text
{{ else }}
Default Text
{{ end }}{{% /if.inline %}}
...`
	output := parseIfStatements(input)
	assert.Equal(t, output, expectedOutput)
}

func TestIfWithHTMLInner(t *testing.T) {
	input := `{% if site.audience == "Cinemo" %} <b>HELLO</b> {% endif %}`
	expectedOutput := `{{< if.inline >}}{{ if (eq $.Page.Site.Params.audience "Cinemo") }} <b>HELLO</b> {{ end }}{{< /if.inline >}}`
	output := parseIfStatements(input)
	assert.Equal(t, output, expectedOutput)
}

func TestIfWithHTMLInner__Shortcode(t *testing.T) {
	input := `{% if site.audience == "Cinemo" %} {{% shortcode %}} <b>HELLO</b> {% endif %}`
	expectedOutput := `{{% when "audience" "Cinemo" %}} {{% shortcode %}} <b>HELLO</b> {{% /when %}}`
	output := parseIfStatements(input)
	assert.Equal(t, output, expectedOutput)
}

func TestStripTooltip(t *testing.T) {
	input := `{{< if.inline >}}{{ if (eq $.Page.Site.Params.audience "Cinemo") }} {{< tooltip "Glossary Item" >}} {{ end }}{{< /if.inline >}}`
	expectedOutput := `{{< if.inline >}}{{ if (eq $.Page.Site.Params.audience "Cinemo") }} Glossary Item {{ end }}{{< /if.inline >}}`
	output := stripTooltips(input)
	assert.Equal(t, output, expectedOutput)
}

func TestStripTooltips(t *testing.T) {
	input := `{{< if.inline >}}{{ if (eq $.Page.Site.Params.audience "Cinemo") }} {{< tooltip "Glossary Item" >}} {{ else }} {{< tooltip "Glossary Item 2" >}} {{ end }}{{< /if.inline >}}`
	expectedOutput := `{{< if.inline >}}{{ if (eq $.Page.Site.Params.audience "Cinemo") }} Glossary Item {{ else }} Glossary Item 2 {{ end }}{{< /if.inline >}}`
	output := stripTooltips(input)
	assert.Equal(t, output, expectedOutput)
}

func TestReplaceComments(t *testing.T) {
	input := `{% comment %} Comment Body {% endcomment %}`
	expectedOutput := `{{< comment.inline >}}
{{/* Comment Body */}}
{{< /comment.inline >}}`
	output := parseComments(input)
	assert.Equal(t, output, expectedOutput)
}

func TestReplaceCommentsMultiline(t *testing.T) {
	input := `{% comment %}
{% comment %}
Comment Body
More Comment
{% endcomment %}`
	expectedOutput := `{{< comment.inline >}}
{{/*
{% comment %}
Comment Body
More Comment
*/}}
{{< /comment.inline >}}`
	output := parseComments(input)
	assert.Equal(t, output, expectedOutput)
}

