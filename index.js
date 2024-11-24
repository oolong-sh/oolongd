import SpriteText from "//unpkg.com/three-spritetext/dist/three-spritetext.mjs";

fetch("http://localhost:11975/graph")
  .then((res) => res.json())
  .then((data) => {
    const Graph = ForceGraph3D()(document.getElementById("3d-graph"))
      .graphData(data)
      .nodeAutoColorBy("group")
      .nodeThreeObject((node) => {
        const sprite = new SpriteText(node.name);
        sprite.material.depthWrite = false;
        sprite.color = node.color;
        sprite.textHeight = 8;
        return sprite;
      })
      .linkOpacity(0.8) // NOTE: baseline opacity can be adjusted, but keep high
      .linkColor((link) => link.color);
    // Spread nodes a little wider
    Graph.d3Force("charge").strength(-120);
  });
