// @ts-check
// Mikhaz Fire Wails Bindings with Dev RPC fallback

async function callBridge(method, ...args) {
  if (typeof window !== 'undefined' && window['go']?.['main']?.['App']?.[method]) {
    return window['go']['main']['App'][method](...args);
  }
  const res = await fetch('http://127.0.0.1:3999/rpc', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ method, args })
  });
  const data = await res.json();
  if (data.error) throw new Error(data.error);
  return data.result;
}

export function AddServiceAccount(arg1) {
  return callBridge('AddServiceAccount', arg1);
}

export function CheckConflict(arg1, arg2, arg3, arg4, arg5) {
  return callBridge('CheckConflict', arg1, arg2, arg3, arg4, arg5);
}

export function CopyNode(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10) {
  return callBridge('CopyNode', arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10);
}

export function CreateAuthUser(arg1, arg2, arg3) {
  return callBridge('CreateAuthUser', arg1, arg2, arg3);
}

export function DeleteAuthUser(arg1, arg2, arg3) {
  return callBridge('DeleteAuthUser', arg1, arg2, arg3);
}

export function DeleteDocument(arg1, arg2, arg3, arg4, arg5) {
  return callBridge('DeleteDocument', arg1, arg2, arg3, arg4, arg5);
}

export function ExportAuthUsers(arg1, arg2) {
  return callBridge('ExportAuthUsers', arg1, arg2);
}

export function GetAccounts() {
  return callBridge('GetAccounts');
}

export function GetDocument(arg1, arg2, arg3, arg4) {
  return callBridge('GetDocument', arg1, arg2, arg3, arg4);
}

export function ImportAuthUsers(arg1, arg2, arg3) {
  return callBridge('ImportAuthUsers', arg1, arg2, arg3);
}

export function ListAuthUsers(arg1, arg2, arg3, arg4, arg5) {
  return callBridge('ListAuthUsers', arg1, arg2, arg3, arg4, arg5);
}

export function ListDocumentSubcollections(arg1, arg2, arg3, arg4) {
  return callBridge('ListDocumentSubcollections', arg1, arg2, arg3, arg4);
}

export function ListDocuments(arg1, arg2, arg3, arg4, arg5, arg6) {
  return callBridge('ListDocuments', arg1, arg2, arg3, arg4, arg5, arg6);
}

export function ListProjects(arg1) {
  return callBridge('ListProjects', arg1);
}

export function ListRootCollections(arg1, arg2, arg3) {
  return callBridge('ListRootCollections', arg1, arg2, arg3);
}

export function LoginGoogleAccount() {
  return callBridge('LoginGoogleAccount');
}

export function RunJSScript(arg1, arg2, arg3, arg4) {
  return callBridge('RunJSScript', arg1, arg2, arg3, arg4);
}

export function RunSimpleQuery(arg1, arg2, arg3, arg4) {
  return callBridge('RunSimpleQuery', arg1, arg2, arg3, arg4);
}

export function SaveDocument(arg1, arg2, arg3, arg4, arg5, arg6) {
  return callBridge('SaveDocument', arg1, arg2, arg3, arg4, arg5, arg6);
}

export function UpdateAuthUser(arg1, arg2, arg3) {
  return callBridge('UpdateAuthUser', arg1, arg2, arg3);
}
