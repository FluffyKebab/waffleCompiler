const { readFileSync } = require("fs");

const run = async () => {
  const buffer = readFileSync("./main.wasm");
  const module = await WebAssembly.compile(buffer);
  const importObject = {
    console: { log: (arg) => console.log(arg) },
  };
  const instance = await WebAssembly.instantiate(module, importObject);
  console.log(instance.exports.main())
};

run();