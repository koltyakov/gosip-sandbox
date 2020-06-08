$ConfigPath = "./config/private.onprem-wap-adfs.json";

function Main() {
  $Context = Get-Context $ConfigPath;

  # CSOM
  $web = $Context.Web;
  $Context.Load($web);
  $Context.ExecuteQuery();

  Write-Host $web.Title;

  $Context | Format-List;

  # PnP PoSH
  Connect-PnPOnline -Url $Context.Url -CurrentCredentials;
  Set-PnPContext -Context $Context;

  Write-Host Get-PnPHomePage;
}

function Get-Context($ConfigPath) {
  $ConfigJson = Get-Content -Raw -Path $ConfigPath | ConvertFrom-Json;
  $SiteUrl = $ConfigJson.siteUrl;
  $Domain = ([System.Uri]$SiteUrl).Host -replace '^www\.';

  $PnPPath = Split-Path -Path (Get-Module -ListAvailable SharePointPnPPowerShell*)[0].Path;
  Write-Host "Using $PnPPath";

  $SpAuthRead = "$PSScriptRoot\bin\wap-auth.exe -configPath $ConfigPath";
  $Cookies = Invoke-Expression $SpAuthRead | ConvertFrom-Json;

  $CookieContainers = "";
  forEach($Prop in $Cookies.PSObject.Properties) {
    $CookieName = $Prop.Name;
    $CookieValue = $Cookies.$CookieName;
    $CookieContainers += @"
      e.WebRequestExecutor.WebRequest.CookieContainer.Add(new Cookie {
        Domain = "$Domain",
        Name = "$CookieName",
        Value = "$CookieValue"
      });
"@;
  }

  $HandlerClass = @"
    using System;
    using System.Net;
    using Microsoft.SharePoint.Client;

    namespace Gosip {
      public static class Auth {
        public static void Apply(ClientContext ctx) {
          ctx.ExecutingWebRequest += RequestEventHandler;
        }
        private static void RequestEventHandler(object sender, WebRequestEventArgs e) {
          e.WebRequestExecutor.WebRequest.CookieContainer = new CookieContainer();
          $CookieContainers
        }
      }
    }
"@

  Add-Type -TypeDefinition $HandlerClass -ReferencedAssemblies "$PnPPath\Microsoft.SharePoint.Client.dll", "$PnPPath\Microsoft.SharePoint.Client.Runtime.dll";
  Add-Type -Path "$PnPPath\Microsoft.SharePoint.Client.dll";
  Add-Type -Path "$PnPPath\Microsoft.SharePoint.Client.Runtime.dll";

  $Context = New-Object Microsoft.SharePoint.Client.ClientContext($SiteUrl);

  [Gosip.Auth]::Apply($Context);

  $Context;
}

Main;
