const fs = require("fs");
const { wabt } = require("wabt")

require("wabt")().then(wabt => {
    fs.readdirSync("./").filter(fileName => fileName.split(".").at(-1) == "wat").forEach(fileName => {
        const module = wabt.parseWat(fileName, fs.readFileSync(fileName, "utf8"));
        const { buffer } = module.toBinary({});
        fs.writeFileSync(fileName.split(".")[0] + ".wasm",  Buffer.from(buffer));
    });
})

