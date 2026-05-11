"use client";

import { atom } from "jotai";
import type { UserResponse } from "@/api/generated/ninjo";

export const userAtom = atom<UserResponse | null>(null);
export const isAuthenticatedAtom = atom((get) => get(userAtom) !== null);
