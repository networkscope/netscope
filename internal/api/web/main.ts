const API_BASE = "/api/v1";

async function api(path: string, opts?: RequestInit) {
    const res = await fetch(`${API_BASE}${path}`, { ...opts, headers: { "Content-Type": "application/json", ...(opts?.headers || {}) } });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return res.json();
}

document.getElementById("scan")?.addEventListener("click", async () => {
    const target = document.getElementById("target")?.value;
    const status = document.getElementById("status");
    if (!target) {
        status.innerHTML = `<div class="error">Enter a target</div>`;
        return;
    }
    status.innerHTML = `<div>Scanning ${target}...</div>`;
    try {
        const data = await api(`/scan?target=${encodeURIComponent(target)}`, { method: "POST" });
        renderAssets(data.assets);
        renderServices(data.services);
        renderFindings(data.findings);
        renderGraph(data.graph);
        status.innerHTML = `<div>Target: ${data.target}</div>`;
    } catch (e) {
        status.innerHTML = `<div class="error">${e}</div>`;
    }
});

function renderAssets(assets: any[]) {
    const el = document.getElementById("assets");
    if (!assets?.length) { el.innerHTML = ""; return; }
    el.innerHTML = `<h2>Assets (${assets.length})</h2><ul>` + assets.map((a: any) => `<li>${a.ID} (${a.Type})</li>`).join("") + `</ul>`;
}

function renderServices(services: any[]) {
    const el = document.getElementById("services");
    if (!services?.length) { el.innerHTML = ""; return; }
    el.innerHTML = `<h2>Services (${services.length})</h2><ul>` + services.map((s: any) => `<li>${s.AssetID}:${s.Port} ${s.Transport}/${s.Protocol}</li>`).join("") + `</ul>`;
}

function renderFindings(findings: any[]) {
    const el = document.getElementById("findings");
    if (!findings?.length) { el.innerHTML = ""; return; }
    el.innerHTML = `<h2>Findings (${findings.length})</h2><ul>` + findings.map((f: any) => `<li>[${f.Severity}] ${f.Title}</li>`).join("") + `</ul>`;
}

function renderGraph(graph: any) {
    const el = document.getElementById("graph");
    if (!graph?.nodes?.length) { el.innerHTML = ""; return; }
    el.innerHTML = `<h2>Graph</h2><pre>` + graph.nodes.map((n: any) => `${n.ID} [${n.Type}] ${n.Label}`).join("\n") + `</pre>`;
}
