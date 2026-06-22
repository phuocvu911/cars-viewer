const express = require("express");
const expressStatic = require("express-static");
const data = require("./data.json");


const app = express();
const PORT = process.env.PORT || 3000;


// Middleware to parse JSON requests
app.use(express.json());


app.get("/api", (req, res) => {
  res.json({
    models: "/api/models",
    models_by_chassis: "/api/models/chassis/{chassis}",
    models_by_brand: "/api/models/brand/{brand}",
    categories: "/api/categories",
    manufacturers: "/api/manufacturers",
  });
});


// Static files (Used to serve images)
app.use("/api/images", expressStatic("img"));


const carModels = data.carModels
, categories = data.categories
, manufacturers = data.manufacturers;

// Find items by the id. If the dataset is large, 
// then preprocessing would make sense. 
app.get("/api/models/chassis/:chassis", (req, res) => {
  
  const chassis = req.params.chassis;  
  const category_id = categories.find((categories) => categories.name.toLowerCase() === chassis.toLowerCase());
  
  if (!category_id || !chassis) {
    return res.status(404).json({ message: "Chassis not found. "});
  }

  const return_arr = [];

  for (let i = 0; i < carModels.length; i++) {
    var obj = carModels[i]; 

    if (obj.categoryId == category_id.id) {
      return_arr.push(obj);
    }
  }

  if (return_arr.length == 0) {
    return res.status(404).json({ message: "No results for the given criteria. " });
  }

  res.json(return_arr);

});

app.get("/api/models/brand/:brand", (req, res) => {
  
  const brand = req.params.brand;  
  const manufacturer = manufacturers.find((manufacturers) => manufacturers.name.toLowerCase() === brand.toLowerCase());
  
  if (!manufacturer || !brand) {
    return res.status(404).json({ message: "Brand not found. "});
  }

  const return_arr = [];

  for (let i = 0; i < carModels.length; i++) {
    var obj = carModels[i]; 

    if (obj.manufacturerId == manufacturer.id) {
      return_arr.push(obj);
    }
  }

  if (return_arr.length == 0) {
    return res.status(404).json({ message: "No results for the given criteria. " });
  }

  res.json(return_arr);

});


// Car Models Handler
app.get("/api/models", (req, res) => {
  res.json(carModels);
});

app.get("/api/models/:id", (req, res) => {
  const id = parseInt(req.params.id);
  const model = carModels.find((model) => model.id === id);

  if (!model) {
    return res.status(404).json({ message: "Car model not found" });
  }

  res.json(model);
});


// Categories Handler
app.get("/api/categories", (req, res) => {
  res.json(categories);
});

app.get("/api/categories/:id", (req, res) => {
  const id = parseInt(req.params.id);
  const category = categories.find((category) => category.id === id);

  if (!category) {
    return res.status(404).json({ message: "Category not found" });
  }

  res.json(category);
});


// Manufacturers Handler
app.get("/api/manufacturers", (req, res) => {
  res.json(manufacturers);
});

app.get("/api/manufacturers/:id", (req, res) => {
  const id = parseInt(req.params.id);
  const manufacturer = manufacturers.find(
    (manufacturer) => manufacturer.id === id
  );

  if (!manufacturer) {
    return res.status(404).json({ message: "Manufacturer not found" });
  }

  res.json(manufacturer);
});


// Serve
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
