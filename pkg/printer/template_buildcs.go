package printer

const templateBuildCs = `
{{ template "copyright" }}

using UnrealBuildTool;

public class {{ .ModuleName }} : ModuleRules
{
	public {{ .ModuleName }} (ReadOnlyTargetRules Target) : base(Target)
	{
		PCHUsage = PCHUsageMode.UseExplicitOrSharedPCHs;

        {{ if .PublicDependencies }}PublicDependencyModuleNames.AddRange(new string[] {
        {{ range $dp := .PublicDependencies }}  {{ printf "%q" $dp }},
        {{ end }}});{{ end }}
        
        {{ if .PrivateDependencies }}PrivateDependencyModuleNames.AddRange(new string[] {
        {{ range $dp := .PrivateDependencies }} {{ printf "%q" $dp }},
        {{ end }}});{{ end }}
    }
}
`
