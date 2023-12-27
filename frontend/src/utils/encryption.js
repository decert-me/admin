export function Encryption () {

    const key = process.env.REACT_APP_ANSWERS_KEY;

    const xorStrings = (input)=>{
        var output = '';
        for (var i = 0; i < input.length; i++) {
            var c = input.charCodeAt(i);
            var k = key.charCodeAt(i % key.length);
            output += String.fromCharCode(c ^ k);
        }
        return output;
    }
    
    const encode = (data)=>{
        return Buffer.from(xorStrings(data), 'utf8').toString('base64');
    }
    
    const decode = (data)=>{
        return xorStrings(Buffer.from(data, 'base64').toString('utf8'));
    }
    

    return {
        encode,
        decode
    }
}
 