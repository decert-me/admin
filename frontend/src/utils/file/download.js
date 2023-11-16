export const download = (url, name) => {
    const x = new window.XMLHttpRequest();
    x.open('GET', url, true);
    x.responseType = 'blob';
    x.onload = () => {
        const url = window.URL.createObjectURL(x.response);
        const a = document.createElement('a');
        a.href = url;
        a.download = name;
        a.click();
    };
    x.send();
}