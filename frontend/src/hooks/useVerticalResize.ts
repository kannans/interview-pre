"use client";

import { useCallback, useEffect, useRef, useState } from "react";

const STORAGE_KEY = "ipc-practice-editor-height";
const DEFAULT_HEIGHT = 320;
const MIN_HEIGHT = 160;
const MAX_HEIGHT = 720;

function readStoredHeight(): number {
  if (typeof window === "undefined") return DEFAULT_HEIGHT;
  const raw = window.localStorage.getItem(STORAGE_KEY);
  const parsed = raw ? Number.parseInt(raw, 10) : NaN;
  if (Number.isNaN(parsed)) return DEFAULT_HEIGHT;
  return Math.min(MAX_HEIGHT, Math.max(MIN_HEIGHT, parsed));
}

type UseVerticalResizeOptions = {
  disabled?: boolean;
};

export function useVerticalResize({ disabled }: UseVerticalResizeOptions = {}) {
  const [height, setHeight] = useState(DEFAULT_HEIGHT);
  const dragging = useRef(false);
  const startY = useRef(0);
  const startHeight = useRef(DEFAULT_HEIGHT);

  useEffect(() => {
    setHeight(readStoredHeight());
  }, []);

  useEffect(() => {
    return () => {
      document.body.style.cursor = "";
      document.body.style.userSelect = "";
    };
  }, []);

  const onPointerDown = useCallback(
    (event: React.PointerEvent<HTMLDivElement>) => {
      if (disabled) return;
      dragging.current = true;
      startY.current = event.clientY;
      startHeight.current = height;
      event.currentTarget.setPointerCapture(event.pointerId);
      document.body.style.cursor = "row-resize";
      document.body.style.userSelect = "none";
    },
    [disabled, height],
  );

  const onPointerMove = useCallback((event: React.PointerEvent<HTMLDivElement>) => {
    if (!dragging.current) return;
    const next = Math.min(
      MAX_HEIGHT,
      Math.max(MIN_HEIGHT, startHeight.current + (event.clientY - startY.current)),
    );
    setHeight(next);
  }, []);

  const onPointerUp = useCallback((event: React.PointerEvent<HTMLDivElement>) => {
    if (!dragging.current) return;
    dragging.current = false;
    event.currentTarget.releasePointerCapture(event.pointerId);
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
    setHeight((current) => {
      window.localStorage.setItem(STORAGE_KEY, String(current));
      return current;
    });
  }, []);

  return {
    height,
    onPointerDown,
    onPointerMove,
    onPointerUp,
  };
}
