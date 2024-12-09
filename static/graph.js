const API_BASE_URL = "http://localhost:11975";
// TODO: pull colors from css
const noteNodeColor = "#B36B42";
const keywordNodeColor = "#77824A";
const loColor = [61, 52, 44];
const hiColor = [215, 196, 132];

function styleGraphData(graphData) {
  const { nodes, links } = graphData;

  const newNodes = nodes.map((node) => ({
    ...node,
    color: node.group === "note" ? noteNodeColor : keywordNodeColor,
  }));

  const newLinks = links.map((link) => {
    const strength = link.strength / 0.4; // TODO: defer to go for normalization
    const color = [
      Math.round(loColor[0] * (1 - strength)) + hiColor[0] * strength,
      Math.round(loColor[1] * (1 - strength)) + hiColor[1] * strength,
      Math.round(loColor[2] * (1 - strength)) + hiColor[2] * strength,
    ];
    return { ...link, color: `rgb(${color[0]}, ${color[1]}, ${color[2]})` };
  });

  return { nodes: newNodes, links: newLinks };
}

let graphInstance = null;
let mode = "2d";

async function loadGraphData() {
  try {
    const response = await fetch(`${API_BASE_URL}/graph`);
    const data = await response.json();
    return styleGraphData(data);
  } catch (error) {
    console.error("Error loading graph data:", error);
    return { nodes: [], links: [] };
  }
}

// TODO: click handler
function clickHandler(node) {
  if (node.group === "note") {
    console.log(`Clicked on node: ${node.name}`);
  }
}

async function initGraph() {
  const graphContainer = document.getElementById("graph-container");

  if (graphInstance) {
    graphInstance._destructor();
  }

  graphContainer.innerHTML = "";

  const graphData = await loadGraphData();

  if (mode === "2d") {
    graphInstance = ForceGraph()(graphContainer)
      .graphData(graphData)
      .backgroundColor("#24211e")
      .onNodeClick(clickHandler);
  } else if (mode === "3d") {
    graphInstance = ForceGraph3D()(graphContainer)
      .graphData(graphData)
      .backgroundColor("#24211e")
      .onNodeClick(clickHandler);
  }
}

document.getElementById("2d-button").addEventListener("click", () => {
  if (mode !== "2d") {
    mode = "2d";
    document.getElementById("2d-button").classList.add("active");
    document.getElementById("3d-button").classList.remove("active");
    initGraph();
  }
});

document.getElementById("3d-button").addEventListener("click", () => {
  if (mode !== "3d") {
    mode = "3d";
    document.getElementById("3d-button").classList.add("active");
    document.getElementById("2d-button").classList.remove("active");
    initGraph();
  }
});

initGraph();
