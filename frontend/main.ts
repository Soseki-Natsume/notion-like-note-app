interface Note {
    id: string;
    title: string;
    content: string;
}

document.addEventListener("DOMContentLoaded", () => {
    const createBtn = document.getElementById("create-btn") as HTMLButtonElement | null;
    const noteIDInput = document.getElementById("note-id") as HTMLInputElement | null;
    const noteTitleInput = document.getElementById("note-title") as HTMLInputElement | null;
    const noteContentInput = document.getElementById("note-content") as HTMLTextAreaElement | null;

    let timeoutID: number | undefined;

    async function saveNote(note: Note): Promise<void>{
        try {
            const response = await fetch("/update", {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json', 
            }, 
            body: JSON.stringify(note), 
        });
        if (!response.ok) {
            console.error('保存失敗', response.statusText);
            return;
        }
            const data = await response.json();
            console.log('保存成功', data);
            
            const activeBtn = document.querySelector(`.select-note-btn[data-id="${note.id}"]`);
            if (activeBtn) {
                const displayTitle = (data.title || note.title || "").trim();
                activeBtn.textContent = displayTitle !== "" ? displayTitle.title : "無題";
            }
        } catch (error) {
            console.error('接続エラー', error);
        }    
    }
    
    function triggerAutoSave(): void {
        if (timeoutID !== undefined) {
            clearTimeout(timeoutID);
        }

        timeoutID = window.setTimeout(() => {
            const currentID = noteIDInput?.value ?? '';
            if (!currentID) return // メモ: 空文字もfalsy

            const currentTitle = noteTitleInput?.value ?? '';
            const currentContent = noteContentInput?.value ?? '';

            saveNote({
                id: currentID, 
                title: currentTitle, 
                content: currentContent
            });
        }, 1000);      
    }

    noteTitleInput?.addEventListener('input', triggerAutoSave)
    noteContentInput?.addEventListener('input', triggerAutoSave) 


    createBtn?.addEventListener("click", async () => {
        try {
            const response = await fetch("/create", {
                method: "POST"
            });
            if (!response.ok) {
                throw new Error("新規作成失敗");
            }
            
            const note: Note = await response.json();
            if (noteIDInput) noteIDInput.value = note.id;
            if (noteTitleInput) noteTitleInput.value = "";
            if (noteContentInput) noteContentInput.value = "";

            const selectBtn = document.createElement("button");
            selectBtn.className = "select-note-btn";
            selectBtn.setAttribute("data-id", note.id);
            selectBtn.textContent = "無題"

            const deleteBtn = document.createElement("button")
            deleteBtn.className = "delete-btn";
            deleteBtn.setAttribute("data-id", note.id);
            deleteBtn.textContent = "- 削除";

            createBtn.after(selectBtn, deleteBtn);
            console.log("新規作成成功 ID:", note.id);
        } catch (error){
            console.error("新規作成エラー", error);
        }
    });
    
    document.body.addEventListener("click", async (event) => {
        const target = event.target as HTMLElement;

        if (target.classList.contains("select-note-btn")) {
            const id = target.getAttribute("data-id");
            if (!id) return;

            try {
                const response = await fetch(`/read?id=${id}`);
                if (!response.ok) throw new Error("読み込み失敗");

                const note: Note = await response.json();
                if (noteIDInput) noteIDInput.value = note.id;
                if (noteTitleInput) noteTitleInput.value = note.title;
                if (noteContentInput) noteContentInput.value = note.content;

                console.log("読み込み成功 ID:", note.id);
            } catch (error) {
                console.error("読み込みエラー", error);
            }
        }

        if (target.classList.contains("delete-btn")) {
            const id = target.getAttribute("data-id");
            if (!id) return;

            const isConfirmed = window.confirm("削除しますか?");
            if (!isConfirmed) return;

            try {
                const response = await fetch("/delete", {
                    method: "DELETE",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ id })
                });

                if (response.ok) {
                    const selectBtn = document.querySelector(`.select-note-btn[data-id="${id}"]`)
                    selectBtn?.remove();
                    target.remove();

                    if (noteIDInput?.value === id) {
                        noteIDInput.value = "";
                        if (noteTitleInput) noteTitleInput.value = "";
                        if (noteContentInput) noteContentInput.value = "";
                    }
                    console.log("削除成功 ID:", id);
                }
            } catch (error) {
                console.error("削除エラー", error);
            }
        }
    });    
});