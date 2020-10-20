$ConfigPath = "./config/private.onprem-wap-adfs.json";

$ConfigJson = Get-Content -Raw -Path $ConfigPath | ConvertFrom-Json;
$SiteUrl = $ConfigJson.siteUrl;
$Domain = ([System.Uri]$SiteUrl).Host -replace '^www\.';

# $SpAuthRead = "go run ./samples/powershell/main.go -configPath $ConfigPath";
$SpAuthRead = "$PSScriptRoot\bin\wap-auth.exe -configPath $ConfigPath";

$Cookies = Invoke-Expression $SpAuthRead | ConvertFrom-Json;
# For a long running processes assume cookie refreshing or proxying

$Session = New-Object Microsoft.PowerShell.Commands.WebRequestSession;

forEach($Prop in $Cookies.PSObject.Properties)
{
  $Cookie = New-Object System.Net.Cookie;
  $CookieName = $Prop.Name;
  $Cookie.Name = $CookieName;
  $Cookie.Value = $Cookies.$CookieName;
  $Cookie.Domain = $Domain;
  $Session.Cookies.Add($Cookie);
}

$Response = Invoke-WebRequest "$SiteUrl/_api/web?$select=Title" `
  -WebSession $Session `
  -Method "GET" `
  -Headers @{"accept"="application/json;odata=verbose"};

$Data = $Response | ConvertFrom-Json;

Write-Host $Data.d.Title;