package printer

const templateHeader = `
{{ template "copyright" }}
#pragma once

#include "CoreMinimal.h"
#include "Modules/ModuleInterface.h"

class I{{ .ModuleName }}Module : public IModuleInterface
{
public:
    virtual void StartupModule() override;
    virtual void ShutdownModule() override;
};

class F{{ .ModuleName }}Module : public I{{ .ModuleName }}Module
{
public:
    virtual void StartupModule() override;
    virtual void ShutdownModule() override;
    virtual bool IsGameModule() const override { {{if .IsGameModule }} return true; {{ else }} return false; {{ end }} }
};
`
