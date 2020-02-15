# -*- coding: utf-8 -*-
"""
Created on Sat Feb 15 22:19:07 2020

@author: Tanmay Thakur
"""
import tensorflow as tf

from tensorflow.keras.models import Model
from tensorflow.keras.layers import Input, SpatialDropout1D, GRU, LSTM,Conv1D, concatenate, Dense
from tensorflow.keras.layers import GlobalAveragePooling1D, GlobalMaxPooling1D 


def get_model(X_train):
    inp = Input(shape=(X_train.shape[1],X_train.shape[2]))
    x = LSTM(256, activation = 'tanh', recurrent_activation = 'sigmoid', recurrent_dropout = 0, unroll = False, use_bias = True, return_sequences = True)(inp)
    y = GRU(128, activation = 'tanh', recurrent_activation = 'sigmoid', recurrent_dropout = 0, unroll = False, use_bias = True,return_sequences = True)(inp)
    x = concatenate([x,y])
    x = SpatialDropout1D(0.2)(x)
    x = Conv1D(64, kernel_size = 3, padding = "same")(x)
    max_pool = GlobalMaxPooling1D()(x)
    avg_pool = GlobalAveragePooling1D()(x)
    x = concatenate([max_pool,avg_pool])
    x = Dense(1024)(x)
    preds = Dense(X_train.shape[2])(x)
    
    model = Model(inp,preds)
    
    return model