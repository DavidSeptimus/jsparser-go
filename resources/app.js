
"use strict";

const fs = require("fs/promises");
const app = express();

const router = express.Router();
router.get('/', async (req, res) => {
    const f = await fs.readFile("/tmp/somefile.txt");
    res.send(f.toString("utf-8"));
});
app.use(router);

app.listen(3000, async () => {
    console.log(`App listening locally at :3000`);
});
