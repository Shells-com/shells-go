Unicode True
;--------------------------------
;Include Modern UI

  !include "MUI2.nsh"

;--------------------------------
; Start
 
 
  !define MUI_PRODUCT "Shells"
  !define MUI_FILE "Shells"
  !define MUI_VERSION ""
 
 
;---------------------------------
;General
!define MUI_ICON "../../res/icon.ico"

  ;Name and file
  Name "Shells"
  OutFile "ShellsSetup.exe"
  ShowInstDetails "nevershow"
  ShowUninstDetails "nevershow"
  

  ;Default installation folder
  InstallDir "$PROGRAMFILES64\Shells\Client"
  
  ;Get installation folder from registry if available
  InstallDirRegKey HKCU "Software\Shells\Client" ""

  RequestExecutionLevel admin

;--------------------------------
;Interface Settings

  !define MUI_WELCOMEPAGE  
  !define MUI_LICENSEPAGE
  !define MUI_DIRECTORYPAGE
  !define MUI_ABORTWARNING
  !define MUI_UNINSTALLER
  !define MUI_UNCONFIRMPAGE
  !define MUI_FINISHPAGE_NOAUTOCLOSE
  !define MUI_FINISHPAGE_RUN
  !define MUI_FINISHPAGE_RUN_TEXT "Launch Shells"
  !define MUI_FINISHPAGE_RUN_FUNCTION "LaunchShells"  

;--------------------------------
;Pages

  !insertmacro MUI_PAGE_LICENSE "license.rtf"
  !insertmacro MUI_PAGE_DIRECTORY
  !insertmacro MUI_PAGE_INSTFILES
  
  !insertmacro MUI_UNPAGE_CONFIRM
  !insertmacro MUI_UNPAGE_INSTFILES
  !insertmacro MUI_PAGE_FINISH  
;--------------------------------
;Languages
 
  !insertmacro MUI_LANGUAGE "English"

;--------------------------------
;Installer Sections

Section "Shells" Installation
  SectionIn RO
  SetOutPath "$INSTDIR"
  
  File "${MUI_FILE}.exe"
  
  ;create desktop shortcut
  CreateShortCut "$DESKTOP\${MUI_PRODUCT}.lnk" "$INSTDIR\${MUI_FILE}.exe" ""
 
  ;create start-menu items
  CreateDirectory "$SMPROGRAMS\${MUI_PRODUCT}"
  CreateShortCut "$SMPROGRAMS\${MUI_PRODUCT}\Uninstall.lnk" "$INSTDIR\Uninstall.exe" "" "$INSTDIR\Uninstall.exe" 0
  CreateShortCut "$SMPROGRAMS\${MUI_PRODUCT}\${MUI_PRODUCT}.lnk" "$INSTDIR\${MUI_FILE}.exe" "" "$INSTDIR\${MUI_FILE}.exe" 0
 
  
  ;Store installation folder
  WriteRegStr HKCU "Software\Shells\Client" "" $INSTDIR
  
;write uninstall information to the registry
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${MUI_PRODUCT}" "DisplayName" "${MUI_PRODUCT} (remove only)"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${MUI_PRODUCT}" "UninstallString" "$INSTDIR\Uninstall.exe"
 
  WriteUninstaller "$INSTDIR\Uninstall.exe"

SectionEnd

;--------------------------------
;Uninstaller Section

Section "Uninstall"

  ;ADD YOUR OWN FILES HERE...

  Delete "$INSTDIR\Uninstall.exe"

  ;Delete Files 
  RMDir /r "$INSTDIR\*.*"    
 
  ;Remove the installation directory
  RMDir "$INSTDIR"
  
  ;Delete Start Menu Shortcuts
  Delete "$DESKTOP\${MUI_PRODUCT}.lnk"
  Delete "$SMPROGRAMS\${MUI_PRODUCT}\*.*"
  RMDir  "$SMPROGRAMS\${MUI_PRODUCT}"
 
  ;Delete Uninstaller And Unistall Registry Entries
  DeleteRegKey HKEY_LOCAL_MACHINE "SOFTWARE\${MUI_PRODUCT}"
  DeleteRegKey HKEY_LOCAL_MACHINE "SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\${MUI_PRODUCT}"  

  DeleteRegKey /ifempty HKCU "Software\Shells\Client"

SectionEnd

Function LaunchShells
  ExecShell "" "$INSTDIR\${MUI_PRODUCT}.exe"
FunctionEnd
