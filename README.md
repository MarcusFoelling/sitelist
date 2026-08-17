[🇩🇪 Deutsch](README.de.md) | 🇬🇧 English

# CLI Tool for managing Microsoft Edge Enterprise Site List (SiteList.xml) written in Go 
![Demo](./demo.gif)

Microsoft Edge can be configured to open certain websites in `Internet Explorer Mode` aka `Compatibility Mode`.
This tool creates and updates the required XML file while ensuring the correct syntax.

Example:
```XML
<site-list version="2">
  <created-by>
    <tool>Manual</tool>
    <version>1</version>
    <date-created>20260812.000000</date-created>
  </created-by>
  <site url="host.domain.tld">
    <compat-mode>IE11</compat-mode>
    <open-in>IE11</open-in>
  </site>
</site-list>
```

Deploy Microsoft Edge Policies (via Group Policy Object) or adjust the following registry keys to enable Enterprise Site List for Internet Explorer Compatibility.
Settings can be made via HKLM or HKCU; HKLM takes precedence over HKCU. Adjust the path to your **SiteList.xml**. It can be located on a local drive, a network share, or web server.

```powershell
$path = "HKCU:\SOFTWARE\Policies\Microsoft\Edge"
New-Item -Path $path -Force
New-ItemProperty -Path $path -Name "InternetExplorerIntegrationLevel" -PropertyType DWORD -Value 1 -Force
New-ItemProperty -Path $path -Name "InternetExplorerIntegrationSiteList" -PropertyType String -Value "C:\SiteList.xml" -Force
```

Use the following URI to check the content and settings of the currently loaded site list file:
```URI
edge://policy/
```

Check with following URI the content and settings of currently loaded sitelist file:
```URI
edge://compat/enterprise
```
