export function getLocalStorage(key) {
  if (key == "all") {
    return localStorage;
  } else {
    return localStorage[key];
  }
}

export function setLocalStorage(key, value) {
  localStorage[key] = value;
  return;
}

export function removeLocalStorage() {
  localStorage.clear();
  return;
}
