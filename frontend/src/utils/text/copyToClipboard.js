import { message } from "antd";


export function copyToClipboard(textToCopy) {
    const tempInput = document.createElement('textarea');
    tempInput.value = textToCopy;
    document.body.appendChild(tempInput);
    tempInput.select();
    document.execCommand('copy');
    document.body.removeChild(tempInput);
    message.success('Copied!')
}