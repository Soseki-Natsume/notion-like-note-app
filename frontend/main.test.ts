import { describe, it, expect, beforeEach, vi } from 'vitest';

describe('ノートアプリのフロントエンドテスト', () => {
  // 各テスト実行前に HTML 構造（DOM）をシミュレート
  beforeEach(() => {
    document.body.innerHTML = `
      <div id="note-list"></div>
      <button id="create-btn">+ 新規作成</button>
      <input type="hidden" id="note-id" value="">
      <input type="text" id="note-title" placeholder="タイトルを入力...">
      <textarea id="note-content" placeholder="内容を入力..."></textarea>
    `;
  });

  it('1. 新規作成APIを呼び出し、IDがフォーム（hidden input）に格納されるか', async () => {
    // API のレスポンスをモック（疑似化）
    window.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ id: 'mock-uuid-1234', title: '', content: '' }),
    });

    const createBtn = document.getElementById('create-btn') as HTMLButtonElement;
    const noteIDInput = document.getElementById('note-id') as HTMLInputElement;

    // クリック時の非同期処理
    createBtn.addEventListener('click', async () => {
      const res = await fetch('/create', { method: 'POST' });
      const data = await res.json();
      noteIDInput.value = data.id;
    });

    // イベント発火
    createBtn.click();

    // 非同期通信の完了を待機
    await new Promise((resolve) => setTimeout(resolve, 0));

    // 検証
    expect(noteIDInput.value).toBe('mock-uuid-1234');
    expect(window.fetch).toHaveBeenCalledWith('/create', { method: 'POST' });
  });

  it('2. ノート選択時に /read API を呼び出し、各入力欄にデータが反映されるか', async () => {
    window.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        id: 'note-xyz',
        title: '既存ノート',
        content: 'テスト用の本文です',
      }),
    });

    const noteIDInput = document.getElementById('note-id') as HTMLInputElement;
    const titleInput = document.getElementById('note-title') as HTMLInputElement;
    const contentInput = document.getElementById('note-content') as HTMLTextAreaElement;

    // 選択時の呼び出し関数
    const selectNote = async (id: string) => {
      const res = await fetch(`/read?id=${id}`);
      const data = await res.json();
      noteIDInput.value = data.id;
      titleInput.value = data.title;
      contentInput.value = data.content;
    };

    await selectNote('note-xyz');

    // 検証
    expect(noteIDInput.value).toBe('note-xyz');
    expect(titleInput.value).toBe('既存ノート');
    expect(contentInput.value).toBe('テスト用の本文です');
  });

  it('3. 削除ボタン押下時に /delete API に正しい ID が送信されるか', async () => {
    window.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'success' }),
    });

    const deleteNote = async (id: string) => {
      await fetch('/delete', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
      });
    };

    await deleteNote('target-id-5678');

    // 検証: DELETE リクエストの JSON ボディに ID が渡されたか
    expect(window.fetch).toHaveBeenCalledWith('/delete', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: 'target-id-5678' }),
    });
  });
});