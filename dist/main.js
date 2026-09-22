//#region frontend/main.ts
document.addEventListener("DOMContentLoaded", () => {
	let e = document.getElementById("create-btn"), t = document.getElementById("note-id"), n = document.getElementById("note-title"), r = document.getElementById("note-content"), i, a = !1;
	async function o(e) {
		try {
			let t = await fetch("/update", {
				method: "PUT",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(e)
			});
			if (!t.ok) {
				console.error("保存失敗", t.statusText);
				return;
			}
			let n = await t.json();
			console.log("保存成功", n);
			let r = document.querySelector(`.select-note-btn[data-id="${e.id}"]`);
			if (r) {
				let t = (n.title || e.title || "").trim();
				r.textContent = t === "" ? "無題" : t;
			}
		} catch (e) {
			console.error("接続エラー", e);
		}
	}
	async function s() {
		try {
			let e = await fetch("/create", { method: "POST" });
			if (!e.ok) throw Error("新規作成失敗");
			let t = await e.json(), n = document.createElement("div");
			n.className = "note-item", n.setAttribute("data-id", t.id);
			let r = document.createElement("button");
			r.className = "select-note-btn", r.setAttribute("data-id", t.id);
			let i = document.createElement("span");
			i.className = "note-title-text", i.textContent = t.title || "無題", r.appendChild(i);
			let a = document.createElement("button");
			a.className = "delete-btn", a.setAttribute("data-id", t.id), a.textContent = "x", n.appendChild(r), n.appendChild(a);
			let o = document.getElementById("note-list");
			return o && o.prepend(n), t;
		} catch (e) {
			return console.error("新規作成エラー", e), null;
		}
	}
	function c() {
		i !== void 0 && clearTimeout(i), i = window.setTimeout(async () => {
			let e = t?.value ?? "";
			if (!e) {
				if (a) return;
				a = !0;
				let n = await s();
				if (a = !1, !n) return;
				e = n.id, t && (t.value = e);
			}
			let i = n?.value ?? "", c = r?.value ?? "";
			o({
				id: e,
				title: i,
				content: c
			});
		}, 1e3);
	}
	n?.addEventListener("input", c), r?.addEventListener("input", c), e?.addEventListener("click", async () => {
		t && (t.value = ""), n && (n.value = ""), r && (r.value = "");
		let e = await s();
		e && t && (t.value = e.id);
	}), document.body.addEventListener("click", async (e) => {
		let i = e.target, a = i?.closest(".select-note-btn");
		if (a) {
			let e = a.getAttribute("data-id");
			if (!e) return;
			try {
				let i = await fetch(`/read?id=${e}`);
				if (!i.ok) throw Error("読み込み失敗");
				let a = await i.json();
				t && (t.value = a.id), n && (n.value = a.title), r && (r.value = a.content), console.log("読み込み成功 ID:", a.id);
			} catch (e) {
				console.error("読み込みエラー", e);
			}
		}
		if (i.classList.contains("delete-btn")) {
			let e = i.getAttribute("data-id");
			if (!e || !window.confirm("削除しますか?")) return;
			try {
				(await fetch("/delete", {
					method: "DELETE",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ id: e })
				})).ok && (document.querySelector(`.select-note-btn[data-id="${e}"]`)?.remove(), i.remove(), t?.value === e && (t.value = "", n && (n.value = ""), r && (r.value = "")), console.log("削除成功 ID:", e));
			} catch (e) {
				console.error("削除エラー", e);
			}
		}
	});
});
//#endregion
