[🇩🇪 Deutsch](README.de.md) | 🇬🇧 English

# Go SiteList Editor
A CLI tool for managing site entries in an XML file.
![Demo](./demo.gif)

MS Edge (Windows Server 2022) kann angewiesen werden bestimmte Sites mittels `Internet Explorer Modus` zu öffnen. Dieses Feature ist bis mindestens 2029 verfügbar.

Kann über GPO oder Registry gesteuert werden, `HKLM` hat Vorrang gegenüber `HKCU`.

sitelist.xml (Name ist beliebig) anlegen und site -Einträge hinzufügen. Ein Eintrag gilt für alle darunter fallenden Domains. Die Datei kann beliebig platziert werden (lokal, Share, UNC, URI). 
```XML
<site-list version="2">
  <created-by>
    <tool>Manual</tool>
    <version>1</version>
    <date-created>20260812.000000</date-created>
  </created-by>
  <site url="nt.amprion.lan">
    <compat-mode>IE11</compat-mode>
    <open-in>IE11</open-in>
  </site>
</site-list>
```

je nach Bedarf unterhalb HKLM oder HKCU Keys hinzufügen / entfernen
```powershell
Remove-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Edge" -Name "InternetExplorerIntegrationLevel"
```

```powershell
$path = "HKCU:\SOFTWARE\Policies\Microsoft\Edge"
New-Item -Path $path -Force
New-ItemProperty -Path $path -Name "InternetExplorerIntegrationLevel" -PropertyType DWORD -Value 1 -Force
New-ItemProperty -Path $path -Name "InternetExplorerIntegrationSiteList" -PropertyType String -Value "P:\sitelist.xml" -Force
```

prüfen über
```URI
edge://policy/
```

```URI
edge://compat/enterprise
```
