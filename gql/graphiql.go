package gql

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

// GraphiQL is an in-browser IDE for exploring GraphiQL APIs.
// This handler returns GraphiQL when requested.
//
// For more information, see https://github.com/graphql/graphiql.

func respond(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}

func errorJSON(msg string) []byte {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `{"error": "%s"}`, msg)
	return buf.Bytes()
}

func PlaygroundHandler(apiURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respond(w, errorJSON("only GET requests are supported"), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err := page.Execute(w, map[string]interface{}{
			"apiURL": apiURL,
		})
		if err != nil {
			panic(err)
		}
	}
}

// https://github.com/graphql/graphiql/blob/main/examples/graphiql-cdn/index.html
var page = template.Must(template.New("graphiql").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>GraphiQL</title>
    <style>
      body {
        height: 100%;
        margin: 0;
        width: 100%;
        overflow: hidden;
      }

      #graphiql {
        height: 100vh;
      }
    </style>

    <!--
      This GraphiQL example depends on Promise and fetch, which are available in
      modern browsers, but can be "polyfilled" for older browsers.
      GraphiQL itself depends on React DOM.
      If you do not want to rely on a CDN, you can host these files locally or
      include them directly in your favored resource bundler.
    -->
    <script
      src="https://unpkg.com/react@17/umd/react.development.js"
      integrity="sha512-Vf2xGDzpqUOEIKO+X2rgTLWPY+65++WPwCHkX2nFMu9IcstumPsf/uKKRd5prX3wOu8Q0GBylRpsDB26R6ExOg=="
      crossorigin="anonymous"
    ></script>
    <script
      src="https://unpkg.com/react-dom@17/umd/react-dom.development.js"
      integrity="sha512-Wr9OKCTtq1anK0hq5bY3X/AvDI5EflDSAh0mE9gma+4hl+kXdTJPKZ3TwLMBcrgUeoY0s3dq9JjhCQc7vddtFg=="
      crossorigin="anonymous"
    ></script>

    <!--
      These two files can be found in the npm module, however you may wish to
      copy them directly into your environment, or perhaps include them in your
      favored resource bundler.
     -->
    <link rel="stylesheet" href="https://unpkg.com/graphiql/graphiql.min.css" />
  </head>

  <body>
    <div id="graphiql">Loading...</div>
    <script
      src="https://unpkg.com/graphiql/graphiql.min.js"
      type="application/javascript"
    ></script>
    <script>
      // https://github.com/graphql/graphiql/blob/main/packages/graphiql/resources/renderExample.js
      
      // Parse the search string to get url parameters.
      var search = window.location.search;
      var parameters = {};
      search
        .substr(1)
        .split('&')
        .forEach(function (entry) {
          var eq = entry.indexOf('=');
          if (eq >= 0) {
            parameters[decodeURIComponent(entry.slice(0, eq))] = decodeURIComponent(
              entry.slice(eq + 1),
            );
          }
        });
      
      // When the query and variables string is edited, update the URL bar so
      // that it can be easily shared.
      function onEditQuery(newQuery) {
        parameters.query = newQuery;
        updateURL();
      }
      
      function onEditVariables(newVariables) {
        parameters.variables = newVariables;
        updateURL();
      }
      
      function onEditHeaders(newHeaders) {
        parameters.headers = newHeaders;
        updateURL();
      }
      
      function onTabChange(tabsState) {
        const activeTab = tabsState.tabs[tabsState.activeTabIndex];
        parameters.query = activeTab.query;
        parameters.variables = activeTab.variables;
        parameters.headers = activeTab.headers;
        updateURL();
      }
      
      function updateURL() {
        var newSearch =
          '?' +
          Object.keys(parameters)
            .filter(function (key) {
              return Boolean(parameters[key]);
            })
            .map(function (key) {
              return (
                encodeURIComponent(key) + '=' + encodeURIComponent(parameters[key])
              );
            })
            .join('&');
        history.replaceState(null, null, newSearch);
      }
      
      // Render <GraphiQL /> into the body.
      // See the README in the top level of this module to learn more about
      // how you can customize GraphiQL by providing different values or
      // additional child elements.
      ReactDOM.render(
        React.createElement(GraphiQL, {
          fetcher: GraphiQL.createFetcher({
            // subscriptionUrl: 'ws://localhost:8081/subscriptions',
            url: {{.apiURL}}
          }),
          query: parameters.query,
          variables: parameters.variables,
          headers: parameters.headers,
          defaultHeaders: parameters.defaultHeaders,
          onEditQuery: onEditQuery,
          onEditVariables: onEditVariables,
          onEditHeaders: onEditHeaders,
          defaultEditorToolsVisibility: true,
          isHeadersEditorEnabled: true,
          shouldPersistHeaders: true,
          onTabChange,
        }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>
`))
