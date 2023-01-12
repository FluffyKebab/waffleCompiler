const fs = require("fs");

const run = async () => {
    fs.readdirSync("./").filter(fileName => fileName.split(".").at(-1) == "wasm").forEach( async (fileName) => {
        const module = await WebAssembly.compile(fs.readFileSync(fileName));
        const instance = await WebAssembly.instantiate(module, {});
        
        instance.exports.test()
        console.log(`${fileName} memory after test`)
        console.log(new Uint32Array( instance.exports.memory.buffer))
    });
};

run();