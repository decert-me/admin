import { message } from "antd";


export function copyToClipboard(textToCopy) {
    navigator.clipboard.writeText(textToCopy)
    .then(() => message.success('Copied!'))
    .catch((err) => message.error('Failed to copy:', err));
}