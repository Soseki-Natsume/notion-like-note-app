"use strict";
const title = document.querySelector('#note-title');
const textarea = document.querySelector('#note-content');
const noteID = "1"; // variable for test
async function saveNote(note) {
    try {
        const response = await fetch(`/update`, {
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
    }
    catch (error) {
        console.log('接続エラー', error);
    }
}
let timeoutID;
function triggerAutoSave() {
    if (timeoutID !== undefined) {
        clearTimeout(timeoutID);
    }
    timeoutID = window.setTimeout(() => {
        const currentTitle = title?.value ?? '';
        const currentContent = textarea?.value ?? '';
        saveNote({
            id: noteID,
            title: currentTitle,
            content: currentContent
        });
    }, 1000);
}
if (title) {
    title?.addEventListener('input', () => {
        triggerAutoSave();
    });
}
if (textarea) {
    textarea.addEventListener('input', () => {
        triggerAutoSave();
    });
}
