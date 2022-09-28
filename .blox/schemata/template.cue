{
	_schema: {
		name:      "Template"
		namespace: "cueblox.pulumi.com"
	}

	runtimes: "dotnet" | "go" | "java" | "nodejs" | "python" | "yaml"

	#Template: {
		_dataset: {
			plural: "templates"
			supportedExtensions: ["yaml"]
		}

		#runtime: {
			name:    runtimes
			options: _
		}

		name:        "${PROJECT}"
		description: "${DESCRIPTION}"
		runtime:     #runtime | runtimes

		// We can optionally handle Pulumi YAML here,
		// with variables and resources, or leave this definition open.
		...
	}
}
