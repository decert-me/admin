export function splitUrl(url) {
    const strArr = url.split("/");
    return {
        link_type: strArr[3],
        challenge: strArr[4],
    }
}