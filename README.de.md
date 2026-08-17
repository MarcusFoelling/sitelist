🇩🇪 Deutsch | [🇬🇧 English](README.md)

# CLI-Tool (in Go) zur Verwaltung der Microsoft Edge Enterprise Site List (SiteList.xml)
![Demo](./demo.gif)

Microsoft Edge kann so konfiguriert werden, dass bestimmte Websites im `Internet Explorer Mode` (auch bekannt als `Kompatibilitätsmodus`) geöffnet werden. Dieses Tool erstellt und aktualisiert die dafür benötigte XML-Datei und stellt dabei die korrekte Syntax sicher.

Beispiel:
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

Über die Microsoft-Edge-Richtlinien (via Group Policy Object) oder die folgenden Registry-Schlüssel kann die Enterprise Site List für den Internet-Explorer-Kompatibilitätsmodus aktiviert werden.
Einstellungen können über HKLM oder HKCU vorgenommen werden; HKLM hat Vorrang vor HKCU. Der Pfad zur **SiteList.xml** ist entsprechend anzupassen – sie kann auf einem lokalen Laufwerk, einer Netzwerkfreigabe oder einem Webserver liegen.
```powershell
$path = "HKCU:\SOFTWARE\Policies\Microsoft\Edge"
New-Item -Path $path -Force
New-ItemProperty -Path $path -Name "InternetExplorerIntegrationLevel" -PropertyType DWORD -Value 1 -Force
New-ItemProperty -Path $path -Name "InternetExplorerIntegrationSiteList" -PropertyType String -Value "C:\SiteList.xml" -Force
```

Über die folgende Adresse können die aktuell konfigurierten Richtlinien geprüft werden:
```URI
edge://policy/
```

Zur Prüfung der aktuell geladenen Site-List-Datei:
```URI
edge://compat/enterprise
```
